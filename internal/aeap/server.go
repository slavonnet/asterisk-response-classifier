package aeap

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/slavonnet/asterisk-response-classifier/internal/classifier"
)

const aeapVersion = "0.1.0"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	port   int
	http   *http.Server
	cls    classifier.Classifier
	mu     sync.Mutex
	closed bool
}

func NewServer(port int, cls classifier.Classifier) *Server {
	s := &Server{port: port, cls: cls}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleWS)
	s.http = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	return s
}

func (s *Server) ListenAndServe() error {
	return s.http.ListenAndServe()
}

func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.closed {
		s.closed = true
		_ = s.http.Close()
	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade: %v", err)
		return
	}
	defer conn.Close()

	session := newSession(conn, s.cls)
	session.run()
}

type session struct {
	conn     *websocket.Conn
	cls      classifier.Classifier
	audioBuf []byte
	codec    string
	language string
	results  []speechResult
	mu       sync.Mutex
}

type speechResult struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

func newSession(conn *websocket.Conn, cls classifier.Classifier) *session {
	return &session{conn: conn, cls: cls, results: make([]speechResult, 0, 8)}
}

func (s *session) run() {
	for {
		msgType, data, err := s.conn.ReadMessage()
		if err != nil {
			return
		}
		if msgType == websocket.BinaryMessage {
			s.mu.Lock()
			s.audioBuf = append(s.audioBuf, data...)
			s.mu.Unlock()
			continue
		}
		s.handleJSON(data)
	}
}

func (s *session) handleJSON(data []byte) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		log.Printf("json parse: %v", err)
		return
	}

	if _, ok := raw["request"]; ok {
		var req aeapRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		resp := s.handleRequest(&req)
		s.sendJSON(resp)
	}
}

type aeapRequest struct {
	Request string          `json:"request"`
	ID      string          `json:"id"`
	Version string          `json:"version"`
	Codecs  []codecEntry    `json:"codecs"`
	Params  json.RawMessage `json:"params"`
}

type codecEntry struct {
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type aeapResponse struct {
	Response string                 `json:"response"`
	ID       string                 `json:"id"`
	Codecs   []codecEntry           `json:"codecs,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
	ErrorMsg string                 `json:"error_msg,omitempty"`
}

func (s *session) handleRequest(req *aeapRequest) aeapResponse {
	resp := aeapResponse{Response: req.Request, ID: req.ID}

	switch req.Request {
	case "setup", "set":
		if err := s.applySetupSet(req); err != nil {
			resp.ErrorMsg = err.Error()
			return resp
		}
		if req.Codecs != nil && len(req.Codecs) > 0 {
			s.codec = req.Codecs[0].Name
			resp.Codecs = []codecEntry{{Name: s.codec}}
		}
		if s.language != "" {
			resp.Params = map[string]interface{}{"language": s.language}
		}
	case "get":
		params, err := s.handleGet(req)
		if err != nil {
			resp.ErrorMsg = err.Error()
			return resp
		}
		resp.Params = params
	default:
		resp.ErrorMsg = fmt.Sprintf("unsupported request: %s", req.Request)
	}
	return resp
}

func (s *session) applySetupSet(req *aeapRequest) error {
	if req.Codecs != nil && len(req.Codecs) > 0 {
		s.codec = req.Codecs[0].Name
	}
	if req.Params == nil {
		return nil
	}
	var params map[string]interface{}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return err
	}
	if lang, ok := params["language"].(string); ok {
		s.language = lang
	}
	if results, ok := params["results"].([]interface{}); ok {
		_ = results // Asterisk sending results to us (not used)
	}
	return nil
}

func (s *session) handleGet(req *aeapRequest) (map[string]interface{}, error) {
	var paramNames []string
	if req.Params != nil {
		if err := json.Unmarshal(req.Params, &paramNames); err != nil {
			var obj map[string]interface{}
			if err2 := json.Unmarshal(req.Params, &obj); err2 != nil {
				return nil, fmt.Errorf("invalid get params")
			}
		}
	}

	out := map[string]interface{}{}
	for _, p := range paramNames {
		switch p {
		case "codec":
			out["codec"] = s.codec
		case "language":
			out["language"] = s.language
		case "results":
			s.mu.Lock()
			// Finalize audio buffer on results poll (end of SpeechBackground)
			if len(s.audioBuf) > 0 {
				r := s.cls.ProcessAudio(s.audioBuf)
				s.results = append(s.results, speechResult{
					Text:  r.Text,
					Score: r.Score,
				})
				s.audioBuf = s.audioBuf[:0]
			}
			out["results"] = s.results
			s.results = s.results[:0]
			s.mu.Unlock()
		default:
			log.Printf("get: unknown param %q", p)
		}
	}
	return out, nil
}

func (s *session) sendJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	if err := s.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("write: %v", err)
	}
}

// PushResult sends interim/final STT result to Asterisk (speech_to_text protocol).
func (s *session) pushResult(text string, score float64) {
	req := map[string]interface{}{
		"request": "set",
		"id":      newID(),
		"params": map[string]interface{}{
			"results": []speechResult{{Text: text, Score: score}},
		},
	}
	s.sendJSON(req)
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

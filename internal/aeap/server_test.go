package aeap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slavonnet/asterisk-response-classifier/internal/classifier"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

func testClassifier(t *testing.T) classifier.Classifier {
	dir := t.TempDir()
	path := filepath.Join(dir, "ivr.yaml")
	_ = os.WriteFile(path, []byte(`positive: ["да"]
negative: ["нет"]
speech_to_phrase: "tcp://127.0.0.1:1"
`), 0o644)
	return classifier.NewService(config.NewLoader(path))
}

func TestSetupHandshake(t *testing.T) {
	srv := NewServer(0, testClassifier(t))
	ts := httptest.NewServer(http.HandlerFunc(srv.handleWS))
	defer ts.Close()
	conn, _, _ := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(ts.URL, "http"), nil)
	defer conn.Close()
	_ = conn.WriteJSON(map[string]interface{}{
		"request": "setup", "id": "1", "version": aeapVersion,
		"codecs": []map[string]string{{"name": "ulaw"}},
	})
	var resp aeapResponse
	_ = conn.ReadJSON(&resp)
	if resp.Response != "setup" {
		t.Fatal(resp)
	}
}

func TestGetResultsUncertainWithoutSTP(t *testing.T) {
	srv := NewServer(0, testClassifier(t))
	ts := httptest.NewServer(http.HandlerFunc(srv.handleWS))
	defer ts.Close()
	conn, _, _ := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(ts.URL, "http"), nil)
	defer conn.Close()
	_ = conn.WriteJSON(map[string]interface{}{
		"request": "setup", "id": "1", "version": aeapVersion,
		"codecs": []map[string]string{{"name": "ulaw"}},
	})
	var setup aeapResponse
	_ = conn.ReadJSON(&setup)
	_ = conn.WriteMessage(websocket.BinaryMessage, []byte{0xff, 0x7f})
	_ = conn.WriteJSON(map[string]interface{}{
		"request": "get", "id": "2", "params": []string{"results"},
	})
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	var getResp aeapResponse
	_ = conn.ReadJSON(&getResp)
	raw, _ := json.Marshal(getResp.Params["results"])
	var results []speechResult
	_ = json.Unmarshal(raw, &results)
	if len(results) != 1 || results[0].Text != "uncertain" {
		t.Fatalf("results = %+v", results)
	}
}

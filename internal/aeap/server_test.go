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
	"github.com/slavonnet/asterisk-response-classifier/internal/phrase"
)

func testService(t *testing.T) classifier.Classifier {
	t.Helper()
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "sentences.yaml")
	_ = os.WriteFile(cfgPath, []byte(`language: ru
lists:
  yes_word:
    values: [{in: "да"}]
  no_word:
    values: [{in: "нет"}]
intents:
  Positive:
    data: [{sentences: ["{yes_word}"]}]
  Negative:
    data: [{sentences: ["{no_word}"]}]
`), 0o644)
	pe, _ := phrase.New("")
	return classifier.NewService(config.NewLoader(cfgPath), pe)
}

func TestSetupHandshake(t *testing.T) {
	srv := NewServer(0, testService(t))
	ts := httptest.NewServer(http.HandlerFunc(srv.handleWS))
	defer ts.Close()
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(ts.URL, "http"), nil)
	if err != nil {
		t.Fatal(err)
	}
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

func TestGetResults(t *testing.T) {
	srv := NewServer(0, testService(t))
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
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	var getResp aeapResponse
	_ = conn.ReadJSON(&getResp)
	raw, _ := json.Marshal(getResp.Params["results"])
	var results []speechResult
	_ = json.Unmarshal(raw, &results)
	if len(results) != 1 || results[0].Text != "uncertain" {
		t.Fatalf("results = %+v", results)
	}
}

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
	"github.com/slavonnet/asterisk-response-classifier/internal/stt"
)

func testService(t *testing.T) *classifier.Service {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "phrases.yaml")
	if err := os.WriteFile(path, []byte("phrases:\n  positive: [\"да\"]\n  negative: [\"нет\"]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return classifier.NewService(config.NewLoader(path), stt.Noop{})
}

func TestSetupHandshake(t *testing.T) {
	srv := NewServer(0, testService(t))
	ts := httptest.NewServer(http.HandlerFunc(srv.handleWS))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	setup := map[string]interface{}{
		"request": "setup",
		"id":      "test-1",
		"version": aeapVersion,
		"codecs":  []map[string]string{{"name": "ulaw"}},
	}
	if err := conn.WriteJSON(setup); err != nil {
		t.Fatalf("write setup: %v", err)
	}

	var resp aeapResponse
	if err := conn.ReadJSON(&resp); err != nil {
		t.Fatalf("read setup response: %v", err)
	}
	if resp.Response != "setup" || resp.ID != "test-1" {
		t.Fatalf("unexpected setup response: %+v", resp)
	}
}

func TestGetResultsAfterAudio(t *testing.T) {
	srv := NewServer(0, testService(t))
	ts := httptest.NewServer(http.HandlerFunc(srv.handleWS))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	_ = conn.WriteJSON(map[string]interface{}{
		"request": "setup",
		"id":      "1",
		"version": aeapVersion,
		"codecs":  []map[string]string{{"name": "ulaw"}},
	})
	var setupResp aeapResponse
	_ = conn.ReadJSON(&setupResp)

	_ = conn.WriteMessage(websocket.BinaryMessage, []byte{0xff, 0x7f, 0x00})

	getReq := map[string]interface{}{
		"request": "get",
		"id":      "2",
		"params":  []string{"results"},
	}
	if err := conn.WriteJSON(getReq); err != nil {
		t.Fatalf("write get: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	var getResp aeapResponse
	if err := conn.ReadJSON(&getResp); err != nil {
		t.Fatalf("read get response: %v", err)
	}

	raw, _ := json.Marshal(getResp.Params["results"])
	var results []speechResult
	_ = json.Unmarshal(raw, &results)
	if len(results) != 1 || results[0].Text != string(classifier.Uncertain) {
		t.Fatalf("results: %+v", results)
	}
}

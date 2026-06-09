package aeap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slavonnet/asterisk-response-classifier/internal/classifier"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

func TestSetupHandshake(t *testing.T) {
	cls := classifier.NewKeywordClassifier(config.PhraseGroups{
		Positive: []string{"да"},
		Negative: []string{"нет"},
	})

	srv := NewServer(0, cls)
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
		"params":  map[string]string{"language": "ru-RU"},
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
	if len(resp.Codecs) != 1 || resp.Codecs[0].Name != "ulaw" {
		t.Fatalf("codecs: %+v", resp.Codecs)
	}
}

func TestGetResultsAfterAudio(t *testing.T) {
	cls := classifier.NewKeywordClassifier(config.PhraseGroups{})
	srv := NewServer(0, cls)
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

	if err := conn.WriteMessage(websocket.BinaryMessage, []byte{0xff, 0x7f, 0x00}); err != nil {
		t.Fatalf("write audio: %v", err)
	}

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
	if getResp.ErrorMsg != "" {
		t.Fatalf("get error: %s", getResp.ErrorMsg)
	}

	raw, err := json.Marshal(getResp.Params["results"])
	if err != nil {
		t.Fatalf("marshal results: %v", err)
	}
	var results []speechResult
	if err := json.Unmarshal(raw, &results); err != nil {
		t.Fatalf("unmarshal results: %v", err)
	}
	if len(results) != 1 || results[0].Text != string(classifier.Uncertain) {
		t.Fatalf("results: %+v", results)
	}
}

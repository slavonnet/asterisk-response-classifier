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
	"github.com/slavonnet/asterisk-response-classifier/internal/audio"
	"github.com/slavonnet/asterisk-response-classifier/internal/classifier"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
)

func testClassifier(t *testing.T) classifier.Classifier {
	t.Helper()
	dir := t.TempDir()
	pos := make([]int16, 1000)
	for i := range pos {
		pos[i] = 2000
	}
	neg := make([]int16, 3000)
	for i := range neg {
		neg[i] = int16(i % 500)
	}
	_ = os.MkdirAll(filepath.Join(dir, "refs"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, "refs/p.ulaw"), audio.EncodeULaw(pos), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "refs/n.ulaw"), audio.EncodeULaw(neg), 0o644)
	cfgPath := filepath.Join(dir, "references.yaml")
	_ = os.WriteFile(cfgPath, []byte(`references:
  positive: ["refs/p.ulaw"]
  negative: ["refs/n.ulaw"]
`), 0o644)
	return classifier.NewSimilarityClassifier(config.NewLoader(cfgPath))
}

func TestSetupHandshake(t *testing.T) {
	srv := NewServer(0, testClassifier(t))
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
	if err := conn.ReadJSON(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Response != "setup" {
		t.Fatalf("resp = %+v", resp)
	}
}

func TestGetResults(t *testing.T) {
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

	_ = conn.WriteMessage(websocket.BinaryMessage, audio.EncodeULaw([]int16{1000, 1000, 1000}))

	_ = conn.WriteJSON(map[string]interface{}{
		"request": "get", "id": "2", "params": []string{"results"},
	})
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	var getResp aeapResponse
	_ = conn.ReadJSON(&getResp)

	raw, _ := json.Marshal(getResp.Params["results"])
	var results []speechResult
	_ = json.Unmarshal(raw, &results)
	if len(results) != 1 {
		t.Fatalf("results = %+v", results)
	}
}

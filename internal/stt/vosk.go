package stt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// VoskWS sends PCM to a Vosk websocket server (alphacep/vosk-server).
// Pure Go — no CGO. Run vosk-server on Pi separately (see deploy/docker-compose.yml).
type VoskWS struct {
	URL string
}

type voskMessage struct {
	Text    string `json:"text"`
	Partial string `json:"partial"`
}

func (v *VoskWS) Transcribe(pcm []byte, sampleRate int) (string, float64, error) {
	if v.URL == "" {
		return "", 0, fmt.Errorf("vosk url empty")
	}
	if len(pcm) == 0 {
		return "", 0, nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(v.URL, nil)
	if err != nil {
		return "", 0, fmt.Errorf("vosk dial: %w", err)
	}
	defer conn.Close()

	cfg := map[string]interface{}{
		"config": map[string]interface{}{"sample_rate": sampleRate},
	}
	if err := conn.WriteJSON(cfg); err != nil {
		return "", 0, fmt.Errorf("vosk config: %w", err)
	}

	const chunk = 8000 // ~0.5s at 8kHz 16-bit
	for i := 0; i < len(pcm); i += chunk {
		end := i + chunk
		if end > len(pcm) {
			end = len(pcm)
		}
		if err := conn.WriteMessage(websocket.BinaryMessage, pcm[i:end]); err != nil {
			return "", 0, fmt.Errorf("vosk audio: %w", err)
		}
	}

	if err := conn.WriteJSON(map[string]bool{"eof": true}); err != nil {
		return "", 0, fmt.Errorf("vosk eof: %w", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	var final string
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg voskMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		if msg.Text != "" {
			final = msg.Text
		}
	}

	if final == "" {
		return "", 0, nil
	}
	return final, 0.85, nil
}

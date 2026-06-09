package wyoming

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	rate     = 16000
	width    = 2
	channels = 1
)

// Transcribe sends PCM 16kHz mono LE to speech-to-phrase (Wyoming ASR).
func Transcribe(addr string, pcm16k []byte) (string, error) {
	host := strings.TrimPrefix(strings.TrimPrefix(addr, "tcp://"), "tcp:")
	conn, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		return "", fmt.Errorf("stp connect %s: %w", host, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(15 * time.Second))

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	if err := writeEvent(w, "transcribe", map[string]any{"language": "ru"}); err != nil {
		return "", err
	}
	if err := writeEvent(w, "audio-start", map[string]any{
		"rate": rate, "width": width, "channels": channels,
	}); err != nil {
		return "", err
	}

	const chunk = 3200 // 100ms @ 16k
	for i := 0; i < len(pcm16k); i += chunk {
		end := i + chunk
		if end > len(pcm16k) {
			end = len(pcm16k)
		}
		if err := writeAudioChunk(w, pcm16k[i:end]); err != nil {
			return "", err
		}
	}
	if err := writeEvent(w, "audio-stop", map[string]any{}); err != nil {
		return "", err
	}
	if err := w.Flush(); err != nil {
		return "", err
	}

	for {
		ev, err := readEvent(r)
		if err != nil {
			return "", err
		}
		if ev.Type == "transcript" {
			text, _ := ev.Data["text"].(string)
			return strings.TrimSpace(text), nil
		}
	}
}

type event struct {
	Type string
	Data map[string]any
}

func writeEvent(w *bufio.Writer, typ string, data map[string]any) error {
	header := map[string]any{"type": typ, "version": "1.0.0"}
	var dataBytes []byte
	if len(data) > 0 {
		var err error
		dataBytes, err = json.Marshal(data)
		if err != nil {
			return err
		}
		header["data_length"] = len(dataBytes)
	}
	line, _ := json.Marshal(header)
	if _, err := w.Write(append(line, '\n')); err != nil {
		return err
	}
	if len(dataBytes) > 0 {
		if _, err := w.Write(dataBytes); err != nil {
			return err
		}
	}
	return nil
}

func writeAudioChunk(w *bufio.Writer, audio []byte) error {
	header := map[string]any{
		"type":            "audio-chunk",
		"version":         "1.0.0",
		"data_length":     48,
		"payload_length":  len(audio),
	}
	data := map[string]any{"rate": rate, "width": width, "channels": channels}
	dataBytes, _ := json.Marshal(data)
	line, _ := json.Marshal(header)
	if _, err := w.Write(append(line, '\n')); err != nil {
		return err
	}
	if _, err := w.Write(dataBytes); err != nil {
		return err
	}
	if _, err := w.Write(audio); err != nil {
		return err
	}
	return nil
}

func readEvent(r *bufio.Reader) (event, error) {
	line, err := r.ReadBytes('\n')
	if err != nil {
		return event{}, err
	}
	var header struct {
		Type          string `json:"type"`
		DataLength    int    `json:"data_length"`
		PayloadLength int    `json:"payload_length"`
	}
	if err := json.Unmarshal(line, &header); err != nil {
		return event{}, err
	}
	ev := event{Type: header.Type, Data: map[string]any{}}
	if header.DataLength > 0 {
		raw := make([]byte, header.DataLength)
		if _, err := io.ReadFull(r, raw); err != nil {
			return event{}, err
		}
		_ = json.Unmarshal(raw, &ev.Data)
	}
	if header.PayloadLength > 0 {
		if _, err := io.ReadFull(r, make([]byte, header.PayloadLength)); err != nil {
			return event{}, err
		}
	}
	return ev, nil
}

// Upsample8kTo16k doubles 8kHz PCM16 to 16kHz (speech-to-phrase expects 16k).
func Upsample8kTo16k(pcm8 []byte) []byte {
	n := len(pcm8) / 2
	if n == 0 {
		return nil
	}
	out := make([]byte, n*4)
	for i := 0; i < n; i++ {
		s := int16(binary.LittleEndian.Uint16(pcm8[i*2:]))
		binary.LittleEndian.PutUint16(out[i*4:], uint16(s))
		next := s
		if i+1 < n {
			next = int16(binary.LittleEndian.Uint16(pcm8[(i+1)*2:]))
		}
		mid := int16((int32(s) + int32(next)) / 2)
		binary.LittleEndian.PutUint16(out[i*4+2:], uint16(mid))
	}
	return out
}

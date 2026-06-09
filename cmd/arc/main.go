package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/slavonnet/asterisk-response-classifier/internal/aeap"
	"github.com/slavonnet/asterisk-response-classifier/internal/classifier"
	"github.com/slavonnet/asterisk-response-classifier/internal/config"
	"github.com/slavonnet/asterisk-response-classifier/internal/stt"
)

var version = "dev"

func main() {
	port := flag.Int("port", 9099, "AEAP WebSocket listen port")
	cfgPath := flag.String("config", "config/phrases.yaml", "phrases config (hot-reloaded)")
	voskURL := flag.String("vosk-url", "ws://127.0.0.1:2700", "Vosk websocket STT (empty = disabled)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		log.Println("arc", version)
		return
	}

	loader := config.NewLoader(*cfgPath)

	var engine stt.Engine = stt.Noop{}
	if *voskURL != "" {
		engine = &stt.VoskWS{URL: *voskURL}
		log.Printf("STT: vosk %s", *voskURL)
	} else {
		log.Printf("STT: disabled (all answers → uncertain)")
	}

	svc := classifier.NewService(loader, engine)
	srv := aeap.NewServer(*port, svc)

	go func() {
		log.Printf("arc %s listening ws://0.0.0.0:%d", version, *port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
	srv.Close()
}

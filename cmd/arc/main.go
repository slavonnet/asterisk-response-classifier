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
	modelPath := flag.String("model", "", "path to Vosk speech model dir (required for STT)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		log.Println("arc", version)
		return
	}

	loader := config.NewLoader(*cfgPath)

	engine, err := stt.NewEngine(*modelPath)
	if err != nil {
		log.Fatalf("stt: %v", err)
	}
	if *modelPath == "" {
		log.Printf("WARNING: -model not set, all answers will be uncertain")
	} else {
		log.Printf("STT: vosk model %s (in-process)", *modelPath)
	}

	if c, ok := engine.(interface{ Close() }); ok {
		defer c.Close()
	}

	svc := classifier.NewService(loader, engine)
	srv := aeap.NewServer(*port, svc)

	go func() {
		log.Printf("arc %s listening ws://127.0.0.1:%d (Asterisk aeap.conf → this URL)", version, *port)
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

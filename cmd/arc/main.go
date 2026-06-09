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
)

func main() {
	port := flag.Int("port", 9099, "AEAP WebSocket listen port")
	cfgPath := flag.String("config", "config/phrases.yaml", "path to phrases/decision tree config")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	cls := classifier.NewKeywordClassifier(cfg.Phrases)

	srv := aeap.NewServer(*port, cls)
	go func() {
		log.Printf("AEAP speech server listening on ws://0.0.0.0:%d", *port)
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

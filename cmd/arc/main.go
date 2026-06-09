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

var version = "dev"

func main() {
	port := flag.Int("port", 9099, "AEAP port for Asterisk")
	cfgPath := flag.String("config", "config/ivr.yaml", "ivr config")
	flag.Parse()

	loader := config.NewLoader(*cfgPath)
	srv := aeap.NewServer(*port, classifier.NewService(loader))

	go func() {
		log.Printf("arc %s — AEAP bridge to speech-to-phrase (Wyoming)", version)
		log.Printf("listening ws://127.0.0.1:%d", *port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	srv.Close()
}

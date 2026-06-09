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
	"github.com/slavonnet/asterisk-response-classifier/internal/phrase"
)

var version = "dev"

func main() {
	port := flag.Int("port", 9099, "AEAP port")
	cfgPath := flag.String("config", "config/sentences.yaml", "speech-to-phrase style sentences")
	modelPath := flag.String("model", "model", "acoustic model dir (in release tarball)")
	flag.Parse()

	loader := config.NewLoader(*cfgPath)
	pe, err := phrase.New(*modelPath)
	if err != nil {
		log.Fatalf("phrase engine: %v", err)
	}
	defer pe.Close()

	srv := aeap.NewServer(*port, classifier.NewService(loader, pe))
	go func() {
		log.Printf("arc %s ws://127.0.0.1:%d (phrase-limited, like speech-to-phrase)", version, *port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	srv.Close()
}

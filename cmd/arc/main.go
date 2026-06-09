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
	port := flag.Int("port", 9099, "AEAP WebSocket port")
	cfgPath := flag.String("config", "config/references.yaml", "reference clips config")
	showVersion := flag.Bool("version", false, "print version")
	flag.Parse()

	if *showVersion {
		log.Println("arc", version)
		return
	}

	loader := config.NewLoader(*cfgPath)
	cls := classifier.NewSimilarityClassifier(loader)
	srv := aeap.NewServer(*port, cls)

	go func() {
		log.Printf("arc %s ws://127.0.0.1:%d (audio similarity, no STT)", version, *port)
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

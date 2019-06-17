package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/smallinsky/talkwithme/config"
)

var (
	configFile = flag.String("config", "config.yaml", "path to server configuration file")
)

func main() {
	flag.Parse()

	cfg, err := config.New(*configFile)
	if err != nil {
		log.Fatalf("[ERR] Failed to get server configuration: %v\n", err)
		os.Exit(-1)
	}

	server, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("[ERROR] Failed to init server: %v", err)

	}
	server.start()
	log.Printf("[INFO] Server is running\n")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	server.close()
}

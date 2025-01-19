package main

import (
	"flag"
	"log"

	"github.com/shashwatrathod/redis-internals/main/config"
	"github.com/shashwatrathod/redis-internals/main/server"
)

// setupFlags initializes the command-line flags for the application.
func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the redis server.")
	flag.IntVar(&config.Port, "port", 7379, "port for the redis server.")
	flag.Parse()
}

// main is the entry point for the application.
func main() {
	setupFlags()
	log.Println("getting started")
	server.RunMultiThreadedSyncTcpServer()
}
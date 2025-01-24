package main

import (
	"flag"
	"log"

	"github.com/shashwatrathod/redis-internals/config"
	"github.com/shashwatrathod/redis-internals/server"
)

// setupFlags initializes the command-line flags for the application.
func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the redis server.")
	flag.IntVar(&config.Port, "port", 7379, "port for the redis server.")
	flag.BoolVar(&config.LogRequest, "log_request", false, "whether to log raw request body.")
	flag.Parse()
}

// main is the entry point for the application.
func main() {
	setupFlags()
	log.Println("getting started")
	server.RunAsyncTcpServer()
}

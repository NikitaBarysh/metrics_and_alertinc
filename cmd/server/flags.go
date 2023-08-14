package main

import (
	"flag"
	"os"
)

var flagRunAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if addr := os.Getenv("ADDRESS"); addr != "" {
		flagRunAddr = addr
	}
}

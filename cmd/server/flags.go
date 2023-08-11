package main

import (
	"flag"
	"os"
)

var url string

func parseFlag() {
	flag.StringVar(&url, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		url = addr
	}
}

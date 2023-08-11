package main

import (
	"flag"
)

var url string

func parseFlag() {
	flag.StringVar(&url, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}

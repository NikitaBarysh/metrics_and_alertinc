package main

import (
	"flag"
	"os"
)

var flagRunAddr string

func parseFlag() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()

	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		flagRunAddr = addr
	}
}

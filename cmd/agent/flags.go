package main

import (
	"flag"
)

type Options struct {
	url            string
	pollInterval   int64
	reportInterval int64
}

var options Options

func parseFlags() {
	flag.StringVar(&options.url, "a", "http://localhost:8080", "server address and port")
	flag.Int64Var(&options.pollInterval, "p", 2, "poll interval")
	flag.Int64Var(&options.reportInterval, "r", 10, "report interval")

	flag.Parse()

}

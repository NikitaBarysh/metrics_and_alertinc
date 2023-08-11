package main

import (
	"flag"
	"os"
	"strconv"
)

type Options struct {
	flagRunAddr    string
	pollInterval   int64
	reportInterval int64
}

var options Options

func parseFlags() {
	flag.StringVar(&options.flagRunAddr, "a", ":8080", "address and port to run server")
	flag.Int64Var(&options.pollInterval, "p", 2, "poll interval")
	flag.Int64Var(&options.reportInterval, "r", 10, "report interval")

	flag.Parse()

	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		options.flagRunAddr = addr
	}

	if interval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			options.reportInterval = value
		}
	}

	if interval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			options.pollInterval = value
		}
	}
}

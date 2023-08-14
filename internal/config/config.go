package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type FlagNames struct {
	flagRunAddr    string
	pollInterval   int64
	reportInterval int64
}

func NewFlagNames() *FlagNames {
	return &FlagNames{}
}

func (f *FlagNames) ParseFlags() {
	flag.StringVar(&f.flagRunAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.Int64Var(&f.pollInterval, "p", 2, "poll interval")
	flag.Int64Var(&f.reportInterval, "r", 10, "report interval")

	flag.Parse()

	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		f.flagRunAddr = addr
	}

	if interval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			f.reportInterval = value
		}
	}

	if interval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			f.pollInterval = value
		}
	}

	if !strings.HasPrefix(f.flagRunAddr, "http") && !strings.HasPrefix(f.flagRunAddr, "https") {
		f.flagRunAddr = "http:/" + f.flagRunAddr
	}
}
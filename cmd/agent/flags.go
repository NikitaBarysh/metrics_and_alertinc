package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type FlagNames struct {
	FlagRunAddr    string
	PollInterval   int64
	ReportInterval int64
}

var flagsName FlagNames

func parseFlags() {
	flag.StringVar(&flagsName.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&flagsName.PollInterval, "p", 2, "poll interval")
	flag.Int64Var(&flagsName.ReportInterval, "r", 10, "report interval")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagsName.FlagRunAddr = envRunAddr
	}

	if interval := os.Getenv("REPORT_INTERVAL"); interval != " " {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			flagsName.ReportInterval = value
		}
	}

	if interval := os.Getenv("POLL_INTERVAL"); interval != "" {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			flagsName.PollInterval = value
		}
	}

	if !strings.HasPrefix(flagsName.FlagRunAddr, "http") &&
		!strings.HasPrefix(flagsName.FlagRunAddr, "https") && !strings.HasPrefix(flagsName.FlagRunAddr, "localhost") {
		flagsName.FlagRunAddr = "http://localhost" + flagsName.FlagRunAddr
	}
}

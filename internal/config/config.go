package config

import (
	"flag"
	"os"
	"strconv"
)

type FlagNames struct {
	FlagRunAddr    string
	PollInterval   int64
	ReportInterval int64
}

func NewFlagNames() *FlagNames {
	return &FlagNames{}
}

func (f *FlagNames) ParseFlags() {
	flag.StringVar(&f.FlagRunAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.Int64Var(&f.PollInterval, "p", 2, "poll interval")
	flag.Int64Var(&f.ReportInterval, "r", 10, "report interval")

	flag.Parse()

	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		f.FlagRunAddr = addr
	}

	if interval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			f.ReportInterval = value
		}
	}

	if interval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			f.PollInterval = value
		}
	}

	//if !strings.HasPrefix(f.FlagRunAddr, "http") && !strings.HasPrefix(f.FlagRunAddr, "https") {
	//	f.FlagRunAddr = "http://" + f.FlagRunAddr
	//}
}

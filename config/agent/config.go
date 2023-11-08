package agent

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	URL            string
	PollInterval   int64
	ReportInterval int64
	Key            string
	Limit          int64
}

func NewAgent() (*Config, error) {
	cfg := new(Config)
	flag.StringVar(&cfg.URL, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&cfg.PollInterval, "p", 2, "poll interval")
	flag.Int64Var(&cfg.ReportInterval, "r", 10, "report interval")
	flag.StringVar(&cfg.Key, "k", "", "sign key")
	flag.Int64Var(&cfg.Limit, "l", 8, "rate limit")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.URL = envRunAddr
	}

	if interval := os.Getenv("REPORT_INTERVAL"); interval != " " {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			cfg.ReportInterval = value
		}
	}

	if interval := os.Getenv("POLL_INTERVAL"); interval != "" {
		if value, err := strconv.ParseInt(interval, 10, 64); err == nil {
			cfg.PollInterval = value
		}
	}

	if key := os.Getenv("KEY"); key != "" {
		cfg.Key = key
	}

	if limit := os.Getenv("RATE_LIMIT"); limit != "" {
		if value, err := strconv.ParseInt(limit, 10, 64); err == nil {
			cfg.Limit = value
		}
	}

	if !strings.HasPrefix(cfg.URL, "http") &&
		!strings.HasPrefix(cfg.URL, "https") && !strings.HasPrefix(cfg.URL, "localhost") {
		cfg.URL = "http://localhost" + cfg.URL
	}

	return cfg, nil
}

package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	URL            string
	PollInterval   int64
	ReportInterval int64
	Key            string
	CryptoKey      string
	Limit          int
	ConfigJSON     string
}

func (m *Config) fromJSON() error {

	data, err := os.ReadFile(m.ConfigJSON)
	if err != nil {
		return fmt.Errorf("cannot read json config: %w", err)
	}

	var settings map[string]interface{}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json settings: %w", err)
	}

	for stype, value := range settings {
		switch stype {
		case "address":
			if m.URL == `` {
				m.URL = value.(string)
			}
		case "report_interval":
			if m.ReportInterval == 0 {
				duration, err := time.ParseDuration(value.(string))
				if err != nil {
					return fmt.Errorf("bad json param 'report_interval': %w", err)
				}
				m.ReportInterval = int64(duration.Seconds())
			}
		case "poll_interval":
			if m.PollInterval == 0 {
				duration, err := time.ParseDuration(value.(string))
				if err != nil {
					return fmt.Errorf("bad json param 'poll_interval': %w", err)
				}
				m.PollInterval = int64(duration.Seconds())
			}
		case "sign_key":
			if m.Key == `` {
				m.Key = value.(string)
			}
		case "crypto_key":
			if m.CryptoKey == `` {
				m.CryptoKey = value.(string)
			}
		}
	}

	return nil
}

func NewAgent() (*Config, error) {
	cfg := new(Config)
	flag.StringVar(&cfg.URL, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&cfg.PollInterval, "p", 2, "poll interval")
	flag.Int64Var(&cfg.ReportInterval, "r", 10, "report interval")
	flag.StringVar(&cfg.Key, "k", "", "sign key")
	flag.IntVar(&cfg.Limit, "l", 8, "rate limit")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "open crypto key")
	flag.StringVar(&cfg.ConfigJSON, "c", "", "json config")
	flag.StringVar(&cfg.ConfigJSON, "config", "", "json config")

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

	cfg.Key = os.Getenv("KEY")

	if limit := os.Getenv("RATE_LIMIT"); limit != "" {
		if value, err := strconv.Atoi(limit); err == nil {
			cfg.Limit = value
		}
	}

	cfg.CryptoKey = os.Getenv("CRYPTO_KEY")

	cfg.ConfigJSON = os.Getenv("CONFIG")

	if cfg.ConfigJSON != "" {
		err := cfg.fromJSON()
		if err != nil {
			return nil, fmt.Errorf("err get config from json: %w", err)
		}
	}

	if !strings.HasPrefix(cfg.URL, "http") &&
		!strings.HasPrefix(cfg.URL, "https") && !strings.HasPrefix(cfg.URL, "localhost") {
		cfg.URL = "http://localhost" + cfg.URL
	}

	return cfg, nil
}

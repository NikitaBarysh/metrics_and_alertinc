package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/encrypt"
)

type Config struct {
	URL            string `json:"address"`
	PollInterval   int64  `json:"poll_interval"`
	ReportInterval int64  `json:"report_interval"`
	Key            string `json:"sign_key"`
	CryptoKey      string `json:"crypto_key"`
	Limit          int
	ConfigJSON     string
	IP             string
	ServiceType    string `json:"service_type"`
}

func (m *Config) fromJSON() error {
	var cfg Config

	data, err := os.ReadFile(m.ConfigJSON)
	if err != nil {
		return fmt.Errorf("cannot read json config: %w", err)
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json settings: %w", err)
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

	if cfg.CryptoKey != `` {
		if err := encrypt.InitEncryptor(cfg.CryptoKey); err != nil {
			log.Fatalf("cannot create encryptor: %s\n", err)
		}
	}

	if !strings.HasPrefix(cfg.URL, "http") &&
		!strings.HasPrefix(cfg.URL, "https") && !strings.HasPrefix(cfg.URL, "localhost") {
		cfg.URL = "http://localhost" + cfg.URL
	}

	conn, err := net.Dial("udp", "127.0.0.1:8080")
	if err != nil {
		return nil, fmt.Errorf("err to connect: %w", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	cfg.IP = localAddr.IP.String()

	return cfg, nil
}

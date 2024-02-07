package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	RunAddr       string
	LogLevel      string
	StoreInterval uint64
	StorePath     string
	CryptoKey     string
	Restore       bool
	DataBaseDSN   string
	Key           string
	ConfigJson    string
}

type Option func(c *Config)

func WithRunAddr(address string) Option {
	return func(c *Config) {
		c.RunAddr = address
	}
}

func LogLevel(level string) Option {
	return func(c *Config) {
		c.LogLevel = level
	}
}

func WithStoreInterval(interval uint64) Option {
	return func(c *Config) {
		c.StoreInterval = interval
	}
}

func WithStorePath(path string) Option {
	return func(c *Config) {
		c.StorePath = path
	}
}

func WithRestore(restore bool) Option {
	return func(c *Config) {
		c.Restore = restore
	}
}

func WithDataBaseDSN(db string) Option {
	return func(c *Config) {
		c.DataBaseDSN = db
	}
}

func WithKey(key string) Option {
	return func(c *Config) {
		c.Key = key
	}
}

func NewConfig(option Environment, options ...Option) (*Config, error) {
	cfg := &Config{
		RunAddr:     option.runAddr,
		LogLevel:    option.logLevel,
		StorePath:   option.storePath,
		DataBaseDSN: option.dataBaseDSN,
		Key:         option.key,
		CryptoKey:   option.cryptoKey,
	}

	restore, err := strconv.ParseBool(option.restore)
	if err != nil {
		return nil, fmt.Errorf("restore error: %w", err)
	}
	cfg.Restore = restore

	duration, err := strconv.ParseInt(option.storeInterval, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("store interval error: %w", err)
	}
	cfg.StoreInterval = uint64(duration)

	if cfg.ConfigJson != "" {
		err := cfg.formJson()
		if err != nil {
			return nil, fmt.Errorf("err to get config: %w", err)
		}
	}

	for _, opt := range options {
		opt(cfg)
	}

	return cfg, nil

}

func (m *Config) formJson() error {

	data, err := os.ReadFile(m.ConfigJson)
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
			if m.RunAddr == `` {
				m.RunAddr = value.(string)
			}
		case "restore":
			if !m.Restore {
				m.Restore = value.(bool)
			}
		case "store_interval":
			if m.StoreInterval == 0 {
				duration, err := time.ParseDuration(value.(string))
				if err != nil {
					return fmt.Errorf("bad json param 'store_interval': %w", err)
				}
				m.StoreInterval = uint64(duration.Seconds())
			}
		case "store_file":
			if m.StorePath == `` {
				m.StorePath = value.(string)
			}
		case "database_dsn":
			if m.DataBaseDSN == `` {
				m.DataBaseDSN = value.(string)
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

type Environment struct {
	runAddr       string
	logLevel      string
	storeInterval string
	storePath     string
	restore       string
	dataBaseDSN   string
	key           string
	cryptoKey     string
	configJson    string
}

func NewServer() (Environment, error) {
	var env Environment
	flag.StringVar(&env.runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&env.logLevel, "l", "info", "log level")
	flag.StringVar(&env.storeInterval, "i", "1", "store interval")
	flag.StringVar(&env.storePath, "f", "/tmp/metrics-db.json", "store path")
	flag.StringVar(&env.restore, "r", "true", "restore")
	flag.StringVar(&env.dataBaseDSN, "d", "", "data base DSN")
	flag.StringVar(&env.key, "k", "", "sign key")
	flag.StringVar(&env.cryptoKey, "crypto-key", "", "private crypto key")
	flag.StringVar(&env.configJson, "c", "", "json config")
	flag.StringVar(&env.configJson, "config", "", "json config")

	flag.Parse()

	if addr := os.Getenv("ADDRESS"); addr != "" {
		env.runAddr = addr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		env.logLevel = envLogLevel
	}

	if storeInterval := os.Getenv("STORE_INTERVAL"); storeInterval != "" {
		env.storeInterval = storeInterval
	}

	if storePath := os.Getenv("FILE_STORAGE_PATH"); storePath != "" {
		env.storePath = storePath
	}

	if restore := os.Getenv("RESTORE"); restore != "" {
		env.restore = restore
	}

	if dataBase := os.Getenv("DATABASE_DSN"); dataBase != "" {
		env.dataBaseDSN = dataBase
	}

	if key, ok := os.LookupEnv("KEY"); ok {
		env.key = key
	}

	if cryptoKey, exist := os.LookupEnv("CRYPTO_KEY"); exist {
		env.cryptoKey = cryptoKey
	}

	if config, ok := os.LookupEnv("CONFIG"); ok {
		env.configJson = config
	}

	return env, nil
}

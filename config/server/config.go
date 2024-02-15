package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RunAddr       string `json:"address"`
	LogLevel      string
	StoreInterval uint64 `json:"store_interval"`
	StorePath     string `json:"store_file"`
	CryptoKey     string `json:"crypto_key"`
	Restore       bool   `json:"restore"`
	DataBaseDSN   string `json:"database_dsn"`
	Key           string `json:"sign_key"`
	ConfigJSON    string
	TrustedSubnet string `json:"trusted_subnet"`
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

	if cfg.ConfigJSON != "" {
		err := cfg.fromJSON()
		if err != nil {
			return nil, fmt.Errorf("err to get config: %w", err)
		}
	}

	for _, opt := range options {
		opt(cfg)
	}

	return cfg, nil

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

type Environment struct {
	runAddr       string
	logLevel      string
	storeInterval string
	storePath     string
	restore       string
	dataBaseDSN   string
	key           string
	cryptoKey     string
	configJSON    string
	trustedSubnet string
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
	flag.StringVar(&env.configJSON, "c", "", "json config")
	flag.StringVar(&env.configJSON, "config", "", "json config")
	flag.StringVar(&env.trustedSubnet, "t", "", "subnet")

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
		env.configJSON = config
	}

	if subnet, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		env.trustedSubnet = subnet
	}

	return env, nil
}

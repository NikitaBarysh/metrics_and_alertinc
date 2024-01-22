package server

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RunAddr       string
	LogLevel      string
	StoreInterval uint64
	StorePath     string
	Restore       bool
	DataBaseDSN   string
	Key           string
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

	for _, opt := range options {
		opt(cfg)
	}

	return cfg, nil

}

type Environment struct {
	runAddr       string
	logLevel      string
	storeInterval string
	storePath     string
	restore       string
	dataBaseDSN   string
	key           string
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

	return env, nil
}

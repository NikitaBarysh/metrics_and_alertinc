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
}

func newConfig(option options) (*Config, error) {
	cfg := &Config{
		RunAddr:   option.runAddr,
		LogLevel:  option.logLevel,
		StorePath: option.storePath,
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

	return cfg, nil
}

type options struct {
	runAddr       string
	logLevel      string
	storeInterval string
	storePath     string
	restore       string
}

func ParseServerConfig() (*Config, error) {
	var option options
	flag.StringVar(&option.runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&option.logLevel, "l", "info", "log level")
	flag.StringVar(&option.storeInterval, "i", "300", "store interval")
	flag.StringVar(&option.storePath, "f", "/tmp/metrics-db.json", "store path")
	flag.StringVar(&option.restore, "r", "true", "restore")

	flag.Parse()

	if addr := os.Getenv("ADDRESS"); addr != "" {
		option.runAddr = addr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		option.logLevel = envLogLevel
	}

	if storeInterval := os.Getenv("STORE_INTERVAL"); storeInterval != "" {
		option.storeInterval = storeInterval
	}

	if storePath := os.Getenv("FILE_STORAGE_PATH"); storePath != "" {
		option.storePath = storePath
	}

	if restore := os.Getenv("RESTORE"); restore != "" {
		option.restore = restore
	}

	return newConfig(option)
}

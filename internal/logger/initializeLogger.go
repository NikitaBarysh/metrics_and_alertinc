package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("ошибка при прасинге уровня логера: %w", err)
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("ошибка при билдинге логера: %w", err)
	}
	Log = zl
	return nil
}

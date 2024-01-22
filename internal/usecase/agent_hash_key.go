// Package usecase - Содержит бизнес логику
package usecase

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/config/agent"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service/hasher"
)

// WithHash - генерируем хэш ключ для агента
func WithHash(cfg *agent.Config) *hasher.Hasher {
	if cfg.Key != "" {
		hash := hasher.NewHasher([]byte(cfg.Key))
		hasher.Sign = hash
		return hash

	}
	return nil
}

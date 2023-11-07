package useCase

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/config/agent"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service/hasher"
)

func WithHash(cfg *agent.Config) *hasher.Hasher {
	if cfg.Key != "" {
		hash := hasher.NewHasher([]byte(cfg.Key))
		hasher.Sign = hash
		//fmt.Println("cfg hash", hash)
		//fmt.Println("cfg Sign", hasher.Sign)
		return hash

	}
	return nil
}

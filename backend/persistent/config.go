package persistent

import (
	"context"
	"sync"

	"github.com/iliyanmotovski/raytracer/backend"
)

type inMemoryConfigRepository struct {
	mu     sync.RWMutex
	config *backend.Config
}

func NewInMemoryConfigRepository() backend.ConfigRepository {
	return &inMemoryConfigRepository{}
}

func (r *inMemoryConfigRepository) Get(ctx context.Context) (*backend.Config, error) {
	if err := checkCtx(ctx); err != nil {
		return &backend.Config{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.config, nil
}

func (r *inMemoryConfigRepository) Upsert(ctx context.Context, cfg *backend.Config) (*backend.Config, error) {
	if err := checkCtx(ctx); err != nil {
		return &backend.Config{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.config = cfg
	return r.config, nil
}

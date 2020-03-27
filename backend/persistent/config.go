package persistent

import (
	"context"
	"sync"

	"github.com/iliyanmotovski/raytracer/backend"
)

// inMemoryConfigRepository is a concrete implementation of the ConfigRepository
// which persists the data in memory, it is concurrent safe
type inMemoryConfigRepository struct {
	mu     sync.RWMutex
	config *backend.Config
}

func NewInMemoryConfigRepository() backend.ConfigRepository {
	return &inMemoryConfigRepository{}
}

// Get is used to get the config from the persistence
func (r *inMemoryConfigRepository) Get(ctx context.Context) (*backend.Config, error) {
	if err := checkCtx(ctx); err != nil {
		return &backend.Config{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.config, nil
}

// Upsert is used for create or update the config
func (r *inMemoryConfigRepository) Upsert(ctx context.Context, cfg *backend.Config) (*backend.Config, error) {
	if err := checkCtx(ctx); err != nil {
		return &backend.Config{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.config = cfg
	return r.config, nil
}

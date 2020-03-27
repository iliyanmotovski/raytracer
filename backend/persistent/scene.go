package persistent

import (
	"context"
	"sync"

	"github.com/iliyanmotovski/raytracer/backend"
)

// inMemorySceneRepository is a concrete implementation of the SceneRepository
// which persists the data in memory, it is concurrent safe
type inMemorySceneRepository struct {
	mu    sync.RWMutex
	scene *backend.Scene
}

func NewInMemorySceneRepository() backend.SceneRepository {
	return &inMemorySceneRepository{}
}

// Get is used to get the scene from the persistence
func (r *inMemorySceneRepository) Get(ctx context.Context) (*backend.Scene, error) {
	if err := checkCtx(ctx); err != nil {
		return &backend.Scene{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.scene, nil
}

// Upsert is used for create or update the scene
func (r *inMemorySceneRepository) Upsert(ctx context.Context, scene *backend.Scene) (*backend.Scene, error) {
	if err := checkCtx(ctx); err != nil {
		return &backend.Scene{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.scene = scene
	return r.scene, nil
}

package persistent_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/persistent"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func TestInMemoryConfigRepository(t *testing.T) {
	config := &backend.Config{
		Light: &vector.Vector{X: 250, Y: 300},
		Scene: &vector.Vector{X: 800, Y: 500},
		Polygons: backend.Polygons{
			{
				VerticesCount: 3,
				Loop: vector.Loop{
					{X: 600, Y: 200},
					{X: 646, Y: 133},
					{X: 646, Y: 261},
				},
			},
		},
	}

	db := persistent.NewInMemoryConfigRepository()
	persisted, err := db.Upsert(context.Background(), config)

	assert.Nil(t, err)
	assert.Equal(t, config, persisted)

	got, err := db.Get(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, got, persisted)
}

func TestInMemoryConfigRepositoryUpsertWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	db := persistent.NewInMemoryConfigRepository()
	_, err := db.Upsert(ctx, &backend.Config{})

	assert.Equal(t, backend.ErrContextCancelled, err)
}

func TestInMemoryConfigRepositoryGetWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	db := persistent.NewInMemoryConfigRepository()
	_, err := db.Get(ctx)

	assert.Equal(t, backend.ErrContextCancelled, err)
}

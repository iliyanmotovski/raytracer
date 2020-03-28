package persistent_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/persistent"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func TestInMemorySceneRepository(t *testing.T) {
	scene := &backend.Scene{
		Width:   800,
		Height:  500,
		LitArea: 60,
		Light:   &vector.Vector{X: 250, Y: 300},
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
		Triangles: backend.Triangles{
			{
				Polygon: backend.Polygon{
					VerticesCount: 3,
					Loop: vector.Loop{
						{X: 600, Y: 200},
						{X: 646, Y: 133},
						{X: 646, Y: 261},
					},
				},
			},
		},
	}

	db := persistent.NewInMemorySceneRepository()
	persisted, err := db.Upsert(context.Background(), scene)

	assert.Nil(t, err)
	assert.Equal(t, scene, persisted)

	got, err := db.Get(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, got, persisted)
}

func TestInMemorySceneRepositoryUpsertWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	db := persistent.NewInMemorySceneRepository()
	_, err := db.Upsert(ctx, &backend.Scene{})

	assert.Equal(t, backend.ErrContextCancelled, err)
}

func TestInMemorySceneRepositoryGetWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	db := persistent.NewInMemorySceneRepository()
	_, err := db.Get(ctx)

	assert.Equal(t, backend.ErrContextCancelled, err)
}

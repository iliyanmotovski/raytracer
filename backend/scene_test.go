package backend_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func TestSceneLoading(t *testing.T) {
	boundariesData := `[{"A":{"X":0,"Y":0},"B":{"X":800,"Y":0}},{"A":{"X":800,"Y":0},"B":{"X":800,"Y":500}},{"A":{"X":800,"Y":500` +
		`},"B":{"X":0,"Y":500}},{"A":{"X":0,"Y":500},"B":{"X":0,"Y":0}},{"A":{"X":600,"Y":200},"B":{"X":646,"Y":133}},{"A":{"X":6` +
		`46,"Y":133},"B":{"X":646,"Y":261}},{"A":{"X":646,"Y":261},"B":{"X":600,"Y":200}}]`

	trianglesData := `[{"Loop":[{"X":250,"Y":300},{"X":0,"Y":0.00010000001532262104},{"X":0.00010000000081777657,"Y":0}]` +
		`,"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0.00010000000081777657,"Y":0},{"X":799.9998999999816,"Y":0}],"Vertic` +
		`esCount":3},{"Loop":[{"X":250,"Y":300},{"X":799.9998999999816,"Y":0},{"X":800,"Y":0.00009999999076204062}],"VerticesCoun` +
		`t":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":0.00009999999076204062},{"X":800,"Y":68.05535809478384}],"VerticesCount":3` +
		`},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":68.05535809478384},{"X":645.9998626101681,"Y":133.00020011127697}],"VerticesCo` +
		`unt":3},{"Loop":[{"X":250,"Y":300},{"X":645.9998626101681,"Y":133.00020011127697},{"X":600.0001098143896,"Y":199.9998400` +
		`5295428}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":600.0001098143896,"Y":199.99984005295428},{"X":600.00007976` +
		`87769,"Y":200.00010578033456}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":600.0000797687769,"Y":200.000105780334` +
		`56},{"X":645.9999228901622,"Y":260.9998977456499}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":645.9999228901622,` +
		`"Y":260.9998977456499},{"X":800,"Y":245.83348590064932}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":245.` +
		`83348590064932},{"X":800,"Y":499.99989999999826}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":499.9998999` +
		`9999826},{"X":799.9998999999918,"Y":500}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":799.9998999999918,"Y":500},` +
		`{"X":0.0000999999892883352,"Y":500}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0.0000999999892883352,"Y":500},{` +
		`"X":0,"Y":499.99990000000764}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0,"Y":499.99990000000764},{"X":0,"Y":0` +
		`.00010000001532262104}],"VerticesCount":3}]`

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

	persisted := backend.NewScene(config)
	persisted.LitArea = 93.38
	json.Unmarshal([]byte(trianglesData), &persisted.Triangles)
	json.Unmarshal([]byte(boundariesData), &persisted.Boundaries)

	sceneRepo := new(backend.FakeSceneRepository)
	sceneRepo.On("Upsert", persisted).Return(persisted, nil)

	scene := backend.NewScene(config)
	loaded, err := scene.Load(context.Background(), sceneRepo)

	assert.Nil(t, err)
	assert.Equal(t, persisted, loaded)

	sceneRepo.AssertExpectations(t)
}

func TestSceneLoadingWithRepositoryFailure(t *testing.T) {
	boundariesData := `[{"A":{"X":0,"Y":0},"B":{"X":800,"Y":0}},{"A":{"X":800,"Y":0},"B":{"X":800,"Y":500}},{"A":{"X":800,"Y":500` +
		`},"B":{"X":0,"Y":500}},{"A":{"X":0,"Y":500},"B":{"X":0,"Y":0}},{"A":{"X":600,"Y":200},"B":{"X":646,"Y":133}},{"A":{"X":6` +
		`46,"Y":133},"B":{"X":646,"Y":261}},{"A":{"X":646,"Y":261},"B":{"X":600,"Y":200}}]`

	trianglesData := `[{"Loop":[{"X":250,"Y":300},{"X":0,"Y":0.00010000001532262104},{"X":0.00010000000081777657,"Y":0}]` +
		`,"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0.00010000000081777657,"Y":0},{"X":799.9998999999816,"Y":0}],"Vertic` +
		`esCount":3},{"Loop":[{"X":250,"Y":300},{"X":799.9998999999816,"Y":0},{"X":800,"Y":0.00009999999076204062}],"VerticesCoun` +
		`t":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":0.00009999999076204062},{"X":800,"Y":68.05535809478384}],"VerticesCount":3` +
		`},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":68.05535809478384},{"X":645.9998626101681,"Y":133.00020011127697}],"VerticesCo` +
		`unt":3},{"Loop":[{"X":250,"Y":300},{"X":645.9998626101681,"Y":133.00020011127697},{"X":600.0001098143896,"Y":199.9998400` +
		`5295428}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":600.0001098143896,"Y":199.99984005295428},{"X":600.00007976` +
		`87769,"Y":200.00010578033456}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":600.0000797687769,"Y":200.000105780334` +
		`56},{"X":645.9999228901622,"Y":260.9998977456499}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":645.9999228901622,` +
		`"Y":260.9998977456499},{"X":800,"Y":245.83348590064932}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":245.` +
		`83348590064932},{"X":800,"Y":499.99989999999826}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":499.9998999` +
		`9999826},{"X":799.9998999999918,"Y":500}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":799.9998999999918,"Y":500},` +
		`{"X":0.0000999999892883352,"Y":500}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0.0000999999892883352,"Y":500},{` +
		`"X":0,"Y":499.99990000000764}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0,"Y":499.99990000000764},{"X":0,"Y":0` +
		`.00010000001532262104}],"VerticesCount":3}]`

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

	persisted := backend.NewScene(config)
	persisted.LitArea = 93.38
	json.Unmarshal([]byte(trianglesData), &persisted.Triangles)
	json.Unmarshal([]byte(boundariesData), &persisted.Boundaries)

	sceneRepo := new(backend.FakeSceneRepository)
	sceneRepo.On("Upsert", persisted).Return(&backend.Scene{}, errors.New("error"))

	scene := backend.NewScene(config)
	_, err := scene.Load(context.Background(), sceneRepo)

	assert.Equal(t, errors.New("error"), err)

	sceneRepo.AssertExpectations(t)
}

func TestSceneReloadDaemon(t *testing.T) {
	boundariesData := `[{"A":{"X":0,"Y":0},"B":{"X":800,"Y":0}},{"A":{"X":800,"Y":0},"B":{"X":800,"Y":500}},{"A":{"X":800,"Y":500` +
		`},"B":{"X":0,"Y":500}},{"A":{"X":0,"Y":500},"B":{"X":0,"Y":0}},{"A":{"X":600,"Y":200},"B":{"X":646,"Y":133}},{"A":{"X":6` +
		`46,"Y":133},"B":{"X":646,"Y":261}},{"A":{"X":646,"Y":261},"B":{"X":600,"Y":200}}]`

	trianglesData := `[{"Loop":[{"X":250,"Y":300},{"X":0,"Y":0.00010000001532262104},{"X":0.00010000000081777657,"Y":0}]` +
		`,"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0.00010000000081777657,"Y":0},{"X":799.9998999999816,"Y":0}],"Vertic` +
		`esCount":3},{"Loop":[{"X":250,"Y":300},{"X":799.9998999999816,"Y":0},{"X":800,"Y":0.00009999999076204062}],"VerticesCoun` +
		`t":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":0.00009999999076204062},{"X":800,"Y":68.05535809478384}],"VerticesCount":3` +
		`},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":68.05535809478384},{"X":645.9998626101681,"Y":133.00020011127697}],"VerticesCo` +
		`unt":3},{"Loop":[{"X":250,"Y":300},{"X":645.9998626101681,"Y":133.00020011127697},{"X":600.0001098143896,"Y":199.9998400` +
		`5295428}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":600.0001098143896,"Y":199.99984005295428},{"X":600.00007976` +
		`87769,"Y":200.00010578033456}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":600.0000797687769,"Y":200.000105780334` +
		`56},{"X":645.9999228901622,"Y":260.9998977456499}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":645.9999228901622,` +
		`"Y":260.9998977456499},{"X":800,"Y":245.83348590064932}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":245.` +
		`83348590064932},{"X":800,"Y":499.99989999999826}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":499.9998999` +
		`9999826},{"X":799.9998999999918,"Y":500}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":799.9998999999918,"Y":500},` +
		`{"X":0.0000999999892883352,"Y":500}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0.0000999999892883352,"Y":500},{` +
		`"X":0,"Y":499.99990000000764}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0,"Y":499.99990000000764},{"X":0,"Y":0` +
		`.00010000001532262104}],"VerticesCount":3}]`

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

	persisted := backend.NewScene(config)
	persisted.LitArea = 93.38
	json.Unmarshal([]byte(trianglesData), &persisted.Triangles)
	json.Unmarshal([]byte(boundariesData), &persisted.Boundaries)

	sceneRepo := new(backend.FakeSceneRepository)
	sceneRepo.On("Upsert", persisted).Return(persisted, nil)

	cc := make(chan *backend.ConfigChan)
	srrc := make(chan *backend.SceneReloadResponse)

	srrcFactory := map[string]chan *backend.SceneReloadResponse{"test": srrc}

	daemon := backend.NewSceneReloadDaemon(sceneRepo, cc, srrcFactory)
	daemon.Start(1)

	cc <- &backend.ConfigChan{
		Ctx:          context.Background(),
		Config:       config,
		ResponseChan: "test",
	}

	loaded := <-srrcFactory["test"]

	assert.Nil(t, loaded.Err)
	assert.Equal(t, persisted, loaded.Scene)

	sceneRepo.AssertExpectations(t)
}

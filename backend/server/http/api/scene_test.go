package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/server/http/api"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func TestGetScene(t *testing.T) {
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

	sceneRepo := new(backend.FakeSceneRepository)
	sceneRepo.On("Get").Return(scene, nil)

	r, _ := http.NewRequest("GET", "/api/v1/scene", nil)
	w := httptest.NewRecorder()

	api.GetScene(sceneRepo).ServeHTTP(w, r)

	wantResponse := `{"Width":800,"Height":500,"LitArea":60,"Light":{"X":250,"Y":300},"Polygons":[[{"X":600,"Y":200},{"X":646,"Y":133},{"X":646,"Y":261}]],"Triangles":[[{"X":600,"Y":200},{"X":646,"Y":133},{"X":646,"Y":261}]]}`

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, wantResponse, strings.TrimSpace(w.Body.String()))

	sceneRepo.AssertExpectations(t)
}

func TestGetSceneWithRepositoryFailure(t *testing.T) {
	sceneRepo := new(backend.FakeSceneRepository)
	sceneRepo.On("Get").Return(&backend.Scene{}, errors.New("error"))

	r, _ := http.NewRequest("GET", "/api/v1/scene", nil)
	w := httptest.NewRecorder()

	api.GetScene(sceneRepo).ServeHTTP(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	sceneRepo.AssertExpectations(t)
}

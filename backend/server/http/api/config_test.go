package api_test

import (
	"bytes"
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

func TestCreateConfiguration(t *testing.T) {
	postData := `{
	"scene": {"x": 800, "y": 500},
	"light": {"x": 250, "y": 300},
	"polygons": [
		[{"x": 600, "y": 200}, {"x": 646, "y": 133}, {"x": 646, "y": 261}]
	]
}`

	scene := &backend.Scene{
		Width:  800,
		Height: 500,
		Light:  &vector.Vector{X: 250, Y: 300},
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

	cc := make(chan *backend.ConfigChan)
	srrc := make(chan *backend.SceneReloadResponse)

	srrcFactory := map[string]chan *backend.SceneReloadResponse{backend.CreateConfigHandler: srrc}

	body := bytes.NewReader([]byte(postData))
	r, _ := http.NewRequest("POST", "/api/v1/scene/config", body)
	w := httptest.NewRecorder()

	go func() {
		wantConfig := &backend.ConfigChan{
			Ctx:          r.Context(),
			Config:       config,
			ResponseChan: backend.CreateConfigHandler,
		}

		gotConfig := <-cc
		assert.Equal(t, wantConfig, gotConfig)

		srrcFactory[backend.CreateConfigHandler] <- &backend.SceneReloadResponse{Err: nil, Scene: scene}
	}()

	api.CreateConfiguration(cc, srrcFactory).ServeHTTP(w, r)

	wantResponse := `{"Light":{"X":250,"Y":300},"Scene":{"X":800,"Y":500},"Polygons":[[{"X":600,"Y":200},{"X":646,"Y":133},{"X":646,"Y":261}]]}`

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, wantResponse, strings.TrimSpace(w.Body.String()))
}

func TestCreateConfigurationFromBrokenJSON(t *testing.T) {
	postData := `broken_json`

	cc := make(chan *backend.ConfigChan)
	srrc := make(chan *backend.SceneReloadResponse)

	srrcFactory := map[string]chan *backend.SceneReloadResponse{backend.CreateConfigHandler: srrc}

	body := bytes.NewReader([]byte(postData))
	r, _ := http.NewRequest("POST", "/api/v1/scene/config", body)
	w := httptest.NewRecorder()

	api.CreateConfiguration(cc, srrcFactory).ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateConfigurationWithSceneReloadResponseError(t *testing.T) {
	postData := `{
	"scene": {"x": 800, "y": 500},
	"light": {"x": 250, "y": 300},
	"polygons": [
		[{"x": 600, "y": 200}, {"x": 646, "y": 133}, {"x": 646, "y": 261}]
	]
}`

	cc := make(chan *backend.ConfigChan)
	srrc := make(chan *backend.SceneReloadResponse)

	srrcFactory := map[string]chan *backend.SceneReloadResponse{backend.CreateConfigHandler: srrc}

	body := bytes.NewReader([]byte(postData))
	r, _ := http.NewRequest("POST", "/api/v1/scene/config", body)
	w := httptest.NewRecorder()

	go func() {
		<-cc
		srrcFactory[backend.CreateConfigHandler] <- &backend.SceneReloadResponse{Err: errors.New("error"), Scene: &backend.Scene{}}
	}()

	api.CreateConfiguration(cc, srrcFactory).ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

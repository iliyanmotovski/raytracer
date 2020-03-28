package api

import (
	"encoding/json"
	"net/http"

	"github.com/iliyanmotovski/raytracer/backend"
)

// GetScene is an http handler which gets the scene from the persistence
// and returns it to the caller
func GetScene(sceneRepo backend.SceneRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scene, err := sceneRepo.Get(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		resp := &sceneDTO{
			Light:     &xy{X: scene.Light.X, Y: scene.Light.Y},
			Width:     scene.Width,
			Height:    scene.Height,
			LitArea:   scene.LitArea,
			Polygons:  scene.Polygons,
			Triangles: scene.Triangles,
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

type sceneDTO struct {
	Width, Height, LitArea float64
	Light                  *xy
	Polygons               backend.Polygons
	Triangles              backend.Triangles
}

func (c *sceneDTO) MarshalJSON() ([]byte, error) {
	dto := &struct {
		Width, Height, LitArea float64
		Light                  *xy
		Polygons               [][]*xy
		Triangles              [][]*xy
	}{}

	dto.Width = c.Width
	dto.Height = c.Height
	dto.LitArea = c.LitArea
	dto.Light = c.Light
	dto.Polygons = make([][]*xy, len(c.Polygons))
	dto.Triangles = make([][]*xy, len(c.Triangles))

	for i, polygon := range c.Polygons {
		poly := make([]*xy, len(polygon.Loop))
		for j, vertice := range polygon.Loop {
			poly[j] = &xy{X: vertice.X, Y: vertice.Y}
		}

		dto.Polygons[i] = poly
	}

	for i, triangle := range c.Triangles {
		tri := make([]*xy, len(triangle.Loop))
		for j, vertice := range triangle.Loop {
			tri[j] = &xy{X: vertice.X, Y: vertice.Y}
		}

		dto.Triangles[i] = tri
	}

	return json.Marshal(dto)
}

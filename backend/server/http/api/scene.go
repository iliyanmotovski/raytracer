package api

import (
	"encoding/json"
	"net/http"

	"github.com/iliyanmotovski/raytracer/backend"
)

func GetScene(sceneRepo backend.SceneRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scene, err := sceneRepo.Get(r.Context())
		if err != nil {
			json.NewEncoder(w).Encode(err)
		}

		resp := &sceneDTO{
			Light:     &XY{X: scene.Light.X, Y: scene.Light.Y},
			Width:     scene.Width,
			Height:    scene.Height,
			LitArea:   scene.LitArea,
			Polygons:  scene.Polygons,
			Triangles: scene.Triangles,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

type sceneDTO struct {
	Width, Height, LitArea float64
	Light                  *XY
	Polygons               backend.Polygons
	Triangles              backend.Triangles
}

func (c *sceneDTO) MarshalJSON() ([]byte, error) {
	dto := &struct {
		Width, Height, LitArea float64
		Light                  *XY
		Polygons               [][]*XY
		Triangles              [][]*XY
	}{}

	dto.Width = c.Width
	dto.Height = c.Height
	dto.LitArea = c.LitArea
	dto.Light = c.Light
	dto.Polygons = make([][]*XY, len(c.Polygons))
	dto.Triangles = make([][]*XY, len(c.Triangles))

	for i, polygon := range c.Polygons {
		poly := make([]*XY, len(polygon.Loop))
		for j, vertice := range polygon.Loop {
			poly[j] = &XY{X: vertice.X, Y: vertice.Y}
		}

		dto.Polygons[i] = poly
	}

	for i, triangle := range c.Triangles {
		tri := make([]*XY, len(triangle.Loop))
		for j, vertice := range triangle.Loop {
			tri[j] = &XY{X: vertice.X, Y: vertice.Y}
		}

		dto.Triangles[i] = tri
	}

	return json.Marshal(dto)
}

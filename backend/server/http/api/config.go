package api

import (
	"encoding/json"
	"net/http"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func CreateConfiguration(cc chan *backend.ConfigChan, srrc backend.SceneReloadResponseChanFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := new(configDTO)

		if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}

		cc <- &backend.ConfigChan{Ctx: r.Context(), Config: dto.adapt(), ResponseChan: backend.CreateConfigHandler}
		created := <-srrc[backend.CreateConfigHandler]
		if created.Err != nil {
			json.NewEncoder(w).Encode(created.Err)
			return
		}

		resp := &configDTO{
			Light:    &XY{X: created.Scene.Light.X, Y: created.Scene.Light.Y},
			Scene:    &XY{X: created.Scene.Width, Y: created.Scene.Height},
			Polygons: created.Scene.Polygons,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

type configDTO struct {
	Light, Scene *XY
	Polygons     backend.Polygons
}

func (c *configDTO) adapt() *backend.Config {
	return &backend.Config{
		Light:    &vector.Vector{X: c.Light.X, Y: c.Light.Y},
		Scene:    &vector.Vector{X: c.Scene.X, Y: c.Scene.Y},
		Polygons: c.Polygons,
	}
}

func (c *configDTO) MarshalJSON() ([]byte, error) {
	dto := &struct {
		Light, Scene *XY
		Polygons     [][]*XY
	}{}

	dto.Light = c.Light
	dto.Scene = c.Scene
	dto.Polygons = make([][]*XY, len(c.Polygons))

	for i, polygon := range c.Polygons {
		poly := make([]*XY, len(polygon.Loop))
		for j, vertice := range polygon.Loop {
			poly[j] = &XY{X: vertice.X, Y: vertice.Y}
		}

		dto.Polygons[i] = poly
	}

	return json.Marshal(dto)
}

func (c *configDTO) UnmarshalJSON(b []byte) error {
	dto := &struct {
		Light, Scene *XY
		Polygons     [][]*XY
	}{}

	if err := json.Unmarshal(b, dto); err != nil {
		return err
	}

	c.Light = dto.Light
	c.Scene = dto.Scene
	c.Polygons = make(backend.Polygons, len(dto.Polygons))

	for i, polygon := range dto.Polygons {
		poly := &backend.Polygon{VerticesCount: len(polygon)}
		for _, vertice := range polygon {
			poly.Loop = append(poly.Loop, &vector.Vector{
				X: vertice.X,
				Y: vertice.Y,
			})
		}

		c.Polygons[i] = poly
	}

	return nil
}

type XY struct {
	X, Y float64
}

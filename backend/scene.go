package backend

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/iliyanmotovski/raytracer/backend/vector"
)

type SceneRepository interface {
	Get(context.Context) (*Scene, error)
	Upsert(context.Context, *Scene) (*Scene, error)
}

type Scene struct {
	Width, Height, LitArea float64
	Light                  *vector.Vector
	Polygons               Polygons
	Triangles              Triangles
	Boundaries             Boundaries
}

func NewScene(width, height float64, light *vector.Vector, polygons Polygons) *Scene {
	b := make(Boundaries, 4)
	b[0] = &Boundary{vector.Edge{A: &vector.Vector{0, 0}, B: &vector.Vector{width, 0}}}
	b[1] = &Boundary{vector.Edge{A: &vector.Vector{width, 0}, B: &vector.Vector{width, height}}}
	b[2] = &Boundary{vector.Edge{A: &vector.Vector{width, height}, B: &vector.Vector{0, height}}}
	b[3] = &Boundary{vector.Edge{A: &vector.Vector{0, height}, B: &vector.Vector{0, 0}}}

	return &Scene{Width: width, Height: height, Light: light, Boundaries: b, Polygons: polygons}
}

func (s *Scene) Load(ctx context.Context, repo SceneRepository) (*Scene, error) {
	triangles, litArea, err := s.Process()
	if err != nil {
		return &Scene{}, err
	}

	log.Printf("lit area is %v percent", litArea)

	scene := &Scene{
		Width:      s.Width,
		Height:     s.Height,
		LitArea:    litArea,
		Light:      s.Light,
		Polygons:   s.Polygons,
		Triangles:  triangles,
		Boundaries: s.Boundaries,
	}

	persisted, err := repo.Upsert(ctx, scene)
	if err != nil {
		return &Scene{}, err
	}

	return persisted, nil
}

func (s *Scene) Process() (Triangles, float64, error) {
	for _, polygon := range s.Polygons {
		s.Boundaries = append(s.Boundaries, polygon.GetBoundaries()...)
	}

	if err := s.Polygons.Validate(s.Width, s.Height); err != nil {
		return Triangles{}, 0, err
	}

	particle := NewParticle(s.Light.X, s.Light.Y, s.Boundaries[0:4])
	triangles := particle.Process(s.Boundaries, s.Polygons)

	totalArea := s.Width * s.Height
	litArea := triangles.Area()

	litAreaPercentage := math.Round(((litArea/totalArea)*100)*100) / 100
	return triangles, litAreaPercentage, nil
}

type SceneReloadDaemon struct {
	sceneRepo               SceneRepository
	configChan              chan *ConfigChan
	sceneReloadResponseChan SceneReloadResponseChanFactory
}

func NewSceneReloadDaemon(sceneRepo SceneRepository, cc chan *ConfigChan, srrc SceneReloadResponseChanFactory) *SceneReloadDaemon {
	return &SceneReloadDaemon{
		sceneRepo:               sceneRepo,
		configChan:              cc,
		sceneReloadResponseChan: srrc,
	}
}

func (d *SceneReloadDaemon) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			for c := range d.configChan {
				start := time.Now()

				scene := NewScene(c.Config.Scene.X, c.Config.Scene.Y, c.Config.Light, c.Config.Polygons)
				loaded, err := scene.Load(c.Ctx, d.sceneRepo)

				srrc := d.sceneReloadResponseChan[c.ResponseChan]
				srrc <- &SceneReloadResponse{
					Scene: loaded,
					Err:   err,
				}

				end := time.Now()
				log.Printf("config containing (%d) polygons processed for: %v", len(c.Config.Polygons), end.Sub(start))
			}
		}()
	}
}

type ConfigChan struct {
	Ctx          context.Context
	Config       *Config
	ResponseChan string
}

type SceneReloadResponse struct {
	Scene *Scene
	Err   error
}

type SceneReloadResponseChanFactory map[string]chan *SceneReloadResponse

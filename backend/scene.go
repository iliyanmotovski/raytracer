package backend

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/iliyanmotovski/raytracer/backend/vector"
)

// SceneRepository is an abstraction over some repository
type SceneRepository interface {
	Get(context.Context) (*Scene, error)
	Upsert(context.Context, *Scene) (*Scene, error)
}

// Scene represents the state of the scene
type Scene struct {
	Width, Height, LitArea float64
	Light                  *vector.Vector
	Polygons               Polygons
	Triangles              Triangles
	Boundaries             Boundaries
}

// NewScene creates a new Scene with the 4 basic boundaries - up, right, down, left in this order
func NewScene(width, height float64, light *vector.Vector, polygons Polygons) *Scene {
	b := make(Boundaries, 4)
	b[0] = &Boundary{vector.Edge{A: &vector.Vector{0, 0}, B: &vector.Vector{width, 0}}}
	b[1] = &Boundary{vector.Edge{A: &vector.Vector{width, 0}, B: &vector.Vector{width, height}}}
	b[2] = &Boundary{vector.Edge{A: &vector.Vector{width, height}, B: &vector.Vector{0, height}}}
	b[3] = &Boundary{vector.Edge{A: &vector.Vector{0, height}, B: &vector.Vector{0, 0}}}

	return &Scene{Width: width, Height: height, Light: light, Boundaries: b, Polygons: polygons}
}

// Load reloads the scene with the new configuration and persists it
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

// Process creates a new particle and casts all rays, returns the triangles
// which represent the lit area, the lit area's area in % of the whole scene
// and an error if any. It Validates the polygons as well
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

// SceneReloadDaemon represents a daemon which listens for new scene configuration
// over the provided channel and generates the new scene according the config,
// then sends back the newly generated scene and an error of any on the provided
// return channel
type SceneReloadDaemon struct {
	sceneRepo               SceneRepository
	configChan              chan *ConfigChan
	sceneReloadResponseChan SceneReloadResponseChanFactory
}

// NewSceneReloadDaemon creates new SceneReloadDaemon
func NewSceneReloadDaemon(sceneRepo SceneRepository, cc chan *ConfigChan, srrc SceneReloadResponseChanFactory) *SceneReloadDaemon {
	return &SceneReloadDaemon{
		sceneRepo:               sceneRepo,
		configChan:              cc,
		sceneReloadResponseChan: srrc,
	}
}

// Start starts the scene reload daemon with provided workers
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

// ConfigChan represents the data sent through the
// actual config chan
type ConfigChan struct {
	Ctx          context.Context
	Config       *Config
	ResponseChan string
}

// SceneReloadResponse represents the response sent back from the daemon
type SceneReloadResponse struct {
	Scene *Scene
	Err   error
}

// SceneReloadResponseChanFactory represents a factory of all return channels
// which the daemon picks with key ResponseChan field from ConfigChan struct
type SceneReloadResponseChanFactory map[string]chan *SceneReloadResponse

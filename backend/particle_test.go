package backend_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func TestNewParticle(t *testing.T) {
	bounds := backend.Boundaries{
		{vector.Edge{A: &vector.Vector{1, 2}, B: &vector.Vector{3, 4}}},
		{vector.Edge{A: &vector.Vector{10, 20}, B: &vector.Vector{30, 40}}},
		{vector.Edge{A: &vector.Vector{100, 200}, B: &vector.Vector{300, 400}}},
		{vector.Edge{A: &vector.Vector{1000, 2000}, B: &vector.Vector{3000, 4000}}},
	}

	got := backend.NewParticle(1000, 2000, bounds)

	want := &backend.Particle{
		Pos: &vector.Vector{1000, 2000},
		Rays: backend.Rays{
			{vector.Edge{&vector.Vector{1000, 2000}, &vector.Vector{0.7071067635088781, 0.7071067988642165}}},
			{vector.Edge{&vector.Vector{1000, 2000}, &vector.Vector{-0.4472135596870563, -0.8944272089063657}}},
			{vector.Edge{&vector.Vector{1000, 2000}, &vector.Vector{-0.4468550357563253, -0.8946063810521436}}},
			{vector.Edge{&vector.Vector{1000, 2000}, &vector.Vector{-0.4472136135691947, -0.8944271819652971}}},
			{vector.Edge{&vector.Vector{1000, 2000}, &vector.Vector{-0.4435516785657893, -0.8962487982929019}}},
			{vector.Edge{&vector.Vector{1000, 2000}, &vector.Vector{-0.44721363525227625, -0.8944271711237557}}},
			{vector.Edge{&vector.Vector{1000, 2000}, &vector.Vector{-0.40081878595899345, -0.9161573559287501}}},
			{vector.Edge{&vector.Vector{1000, 2000}, &vector.Vector{0, -1}}},
		},
	}

	assert.Equal(t, want, got)
}

func TestParticleProcess(t *testing.T) {
	screenBounds := backend.Boundaries{
		{vector.Edge{A: &vector.Vector{0, 0}, B: &vector.Vector{800, 0}}},
		{vector.Edge{A: &vector.Vector{800, 0}, B: &vector.Vector{800, 800}}},
		{vector.Edge{A: &vector.Vector{800, 800}, B: &vector.Vector{0, 800}}},
		{vector.Edge{A: &vector.Vector{0, 800}, B: &vector.Vector{0, 0}}},
	}

	particle := backend.NewParticle(250, 300, screenBounds)

	poly := &backend.Polygon{
		VerticesCount: 3,
		Loop: vector.Loop{
			{X: 600, Y: 200},
			{X: 646, Y: 133},
			{X: 646, Y: 261},
		},
	}

	screenBounds = append(screenBounds, poly.GetBoundaries()...)

	triangles := particle.Process(screenBounds, backend.Polygons{poly})

	got, _ := json.Marshal(triangles)

	want := `[{"Loop":[{"X":250,"Y":300},{"X":0,"Y":0.00010000001532262104},{"X":0.00010000000081777657,"Y":0}],"VerticesCount":3},` +
		`{"Loop":[{"X":250,"Y":300},{"X":0.00010000000081777657,"Y":0},{"X":799.9998999999816,"Y":0}],"VerticesCount":3},{"Loop":[{"X":250,"Y` +
		`":300},{"X":799.9998999999816,"Y":0},{"X":800,"Y":0.00009999999076204062}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y"` +
		`:0.00009999999076204062},{"X":800,"Y":68.05535809478384}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":68.0553580947838` +
		`4},{"X":645.9998626101681,"Y":133.00020011127697}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":645.9998626101681,"Y":133.00020` +
		`011127697},{"X":600.0001098143896,"Y":199.99984005295428}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":600.0001098143896,"Y":19` +
		`9.99984005295428},{"X":600.0000797687769,"Y":200.00010578033456}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":600.0000797687769` +
		`,"Y":200.00010578033456},{"X":645.9999228901622,"Y":260.9998977456499}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":645.9999228` +
		`901622,"Y":260.9998977456499},{"X":800,"Y":245.83348590064932}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":245.8334859` +
		`0064932},{"X":800,"Y":799.9999000000123}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":800,"Y":799.9999000000123},{"X":799.99990` +
		`00000137,"Y":800}],"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":799.9999000000137,"Y":800},{"X":0.00009999999485899025,"Y":800}],` +
		`"VerticesCount":3},{"Loop":[{"X":250,"Y":300},{"X":0.00009999999485899025,"Y":800},{"X":0,"Y":799.9999000000131}],"VerticesCount":3},{` +
		`"Loop":[{"X":250,"Y":300},{"X":0,"Y":799.9999000000131},{"X":0,"Y":0.00010000001532262104}],"VerticesCount":3}]`

	assert.Equal(t, want, string(got))
}

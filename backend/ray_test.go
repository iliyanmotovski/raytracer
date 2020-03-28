package backend_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func TestCastRayAgainstBoundary(t *testing.T) {
	ray := backend.NewRay(&vector.Vector{400, 500})
	ray.SetDir(400, 0)

	boundary := &backend.Boundary{Edge: vector.Edge{A: &vector.Vector{0, 250}, B: &vector.Vector{800, 250}}}

	got, ok := ray.Cast(boundary)

	want := &vector.Vector{400, 250}

	assert.True(t, ok)
	assert.Equal(t, want, got)
}

func TestCastRayAgainstBoundaryNotIntersecting(t *testing.T) {
	ray := backend.NewRay(&vector.Vector{400, 500})
	ray.SetDir(800, 800)

	boundary := &backend.Boundary{Edge: vector.Edge{A: &vector.Vector{0, 250}, B: &vector.Vector{800, 250}}}

	got, ok := ray.Cast(boundary)

	assert.False(t, ok)
	assert.Equal(t, (*vector.Vector)(nil), got)
}

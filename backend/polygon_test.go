package backend_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func TestGetPolygonBoundaries(t *testing.T) {
	poly := &backend.Polygon{
		VerticesCount: 3,
		Loop: vector.Loop{
			{X: 600, Y: 200},
			{X: 646, Y: 133},
			{X: 646, Y: 261},
		},
	}

	got := poly.GetBoundaries()

	want := backend.Boundaries{
		{vector.Edge{A: &vector.Vector{600, 200}, B: &vector.Vector{646, 133}}},
		{vector.Edge{A: &vector.Vector{646, 133}, B: &vector.Vector{646, 261}}},
		{vector.Edge{A: &vector.Vector{646, 261}, B: &vector.Vector{600, 200}}},
	}

	assert.Equal(t, want, got)
}

func TestIsPointContainedInPolygon(t *testing.T) {
	cases := []*struct {
		point *vector.Vector
		want  bool
	}{
		{
			&vector.Vector{601, 201},
			true,
		},
		{
			&vector.Vector{599, 199},
			false,
		},
	}

	poly := &backend.Polygon{
		VerticesCount: 3,
		Loop: vector.Loop{
			{X: 600, Y: 200},
			{X: 646, Y: 133},
			{X: 646, Y: 261},
		},
	}

	for i, c := range cases {
		got := poly.IsPointContainedInPolygon(c.point)
		assert.Equal(t, c.want, got, fmt.Sprintf("case failed: %v", i))
	}
}

func TestPolygonContainsVertice(t *testing.T) {
	cases := []*struct {
		point *vector.Vector
		want  bool
	}{
		{
			&vector.Vector{646, 133},
			true,
		},
		{
			&vector.Vector{10000, 10000},
			false,
		},
	}

	poly := &backend.Polygon{
		VerticesCount: 3,
		Loop: vector.Loop{
			{X: 600, Y: 200},
			{X: 646, Y: 133},
			{X: 646, Y: 261},
		},
	}

	for i, c := range cases {
		got := poly.ContainsVertice(c.point)
		assert.Equal(t, c.want, got, fmt.Sprintf("case failed: %v", i))
	}
}

func TestPolygonIsConvex(t *testing.T) {
	cases := []*struct {
		poly *backend.Polygon
		want bool
	}{
		{
			&backend.Polygon{
				VerticesCount: 3,
				Loop: vector.Loop{
					{X: 600, Y: 200},
					{X: 646, Y: 133},
					{X: 646, Y: 261},
				},
			},
			true,
		},
		{
			&backend.Polygon{
				VerticesCount: 6,
				Loop: vector.Loop{
					{X: 131, Y: 188},
					{X: 54, Y: 136},
					{X: 200, Y: 150},
					{X: 220, Y: 32},
					{X: 238, Y: 114},
					{X: 209, Y: 163},
				},
			},
			false,
		},
	}

	for i, c := range cases {
		got := c.poly.IsConvex()
		assert.Equal(t, c.want, got, fmt.Sprintf("case failed: %v", i))
	}
}

func TestValidatePolygons(t *testing.T) {
	polygons := backend.Polygons{
		{
			VerticesCount: 3,
			Loop: vector.Loop{
				{X: 600, Y: 200},
				{X: 646, Y: 133},
				{X: 646, Y: 261},
			},
		},
		{
			VerticesCount: 6,
			Loop: vector.Loop{
				{X: 131, Y: 188},
				{X: 54, Y: 136},
				{X: 86, Y: 32},
				{X: 220, Y: 32},
				{X: 238, Y: 114},
				{X: 209, Y: 163},
			},
		},
	}

	got := polygons.Validate(800, 500)
	assert.Nil(t, got)
}

func TestValidatePolygonsWithPointOutsideOfTheScene(t *testing.T) {
	polygons := backend.Polygons{
		{
			VerticesCount: 3,
			Loop: vector.Loop{
				{X: 600, Y: 200},
				{X: 850, Y: 550},
				{X: 646, Y: 261},
			},
		},
	}

	got := polygons.Validate(800, 500)
	want := errors.New("point X: 850 , Y: 550 is outside the scene")
	assert.Equal(t, want, got)
}

func TestValidateNonConvexPolygons(t *testing.T) {
	polygons := backend.Polygons{
		{
			VerticesCount: 6,
			Loop: vector.Loop{
				{X: 131, Y: 188},
				{X: 54, Y: 136},
				{X: 200, Y: 150},
				{X: 220, Y: 32},
				{X: 238, Y: 114},
				{X: 209, Y: 163},
			},
		},
	}

	got := polygons.Validate(800, 500)
	want := errors.New("polygon is not convex")
	assert.Equal(t, want, got)
}

func TestValidateOverlappingPolygons(t *testing.T) {
	polygons := backend.Polygons{
		{
			VerticesCount: 3,
			Loop: vector.Loop{
				{X: 600, Y: 200},
				{X: 646, Y: 133},
				{X: 646, Y: 261},
			},
		},
		{
			VerticesCount: 6,
			Loop: vector.Loop{
				{X: 601, Y: 201},
				{X: 54, Y: 136},
				{X: 86, Y: 32},
				{X: 220, Y: 32},
				{X: 238, Y: 114},
				{X: 209, Y: 163},
			},
		},
	}

	got := polygons.Validate(800, 500)
	want := errors.New("point X: 601 , Y: 201 is inside another polygon")
	assert.Equal(t, want, got)
}

func TestGetTriangleArea(t *testing.T) {
	triangle := &backend.Triangle{
		Polygon: backend.Polygon{
			VerticesCount: 3,
			Loop: vector.Loop{
				{X: 600, Y: 200},
				{X: 646, Y: 133},
				{X: 646, Y: 261},
			},
		},
	}

	got := triangle.Area()
	assert.Equal(t, float64(2944), got)
}

func TestGetTrianglesArea(t *testing.T) {
	triangles := backend.Triangles{
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
	}

	got := triangles.Area()
	assert.Equal(t, float64(5888), got)
}

package backend

import (
	"testing"

	"github.com/iliyanmotovski/raytracer/backend/vector"
)

func TestArea(t *testing.T) {
	tr := &Triangle{Polygon{
		Loop: vector.Loop{
			{-4, -5},
			{-4, 13},
			{25, 14},
		},
		VerticesCount: 3,
	}}

	tr = tr
}

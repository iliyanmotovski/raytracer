package backend

import (
	"math"
	"sort"

	"github.com/iliyanmotovski/raytracer/backend/vector"
)

// Particle represents a point from where rays of "light" emit
type Particle struct {
	Pos  *vector.Vector
	Rays Rays
}

// Creates new Particle with given position and sets directory of 8 base rays
// to the 4 corners of the screen, 2 for each corner, with a very small
// offset to the left and right of the center of the corner
func NewParticle(x, y float64, sceneEdgesBounds Boundaries) *Particle {
	baseRays := make(Rays, 8)

	top := sceneEdgesBounds[0]
	right := sceneEdgesBounds[1]
	bottom := sceneEdgesBounds[2]
	left := sceneEdgesBounds[3]

	for i := range baseRays {
		baseRays[i] = NewRay(&vector.Vector{X: x, Y: y})
	}

	baseRays[0].SetDir(left.B.X, left.B.Y+0.0001)
	baseRays[1].SetDir(top.A.X+0.0001, top.A.Y)

	baseRays[2].SetDir(top.B.X-0.0001, top.B.Y)
	baseRays[3].SetDir(right.A.X, right.A.Y+0.0001)

	baseRays[4].SetDir(right.B.X, right.B.Y-0.0001)
	baseRays[5].SetDir(bottom.A.X-0.0001, bottom.A.Y)

	baseRays[6].SetDir(bottom.B.X+0.0001, bottom.B.Y)
	baseRays[7].SetDir(left.A.X, left.A.Y-0.0001)

	return &Particle{Pos: &vector.Vector{X: x, Y: y}, Rays: baseRays}
}

// Process casts the rays
func (p *Particle) Process(boundaries Boundaries, polygons Polygons) Triangles {
	// Adds 2 rays for each polygon vertice and sets their direction with a very
	// small offset to the left and right of the vertice
	p.SetRaysDirToPolyVertices(polygons)
	// sorts the rays clockwise by angle
	p.SortRaysClockwise()

	edges := vector.Loop{}
	for _, ray := range p.Rays {
		var closest *vector.Vector
		lastDistance := math.Inf(1)

		for _, boundary := range boundaries {
			// casts the ray against each boundary
			intersection, ok := ray.Cast(boundary)
			if ok {
				// records the closest point of intersection
				// to the starting point of the ray
				distance := ray.A.Distance(*intersection)
				if distance < lastDistance {
					lastDistance = distance
					closest = intersection
				}
			}
		}

		if closest != nil {
			edges = append(edges, closest)
		}
	}

	return NewClockwiseTriangleFan(p.Pos, edges)
}

// SetRaysDirToPolyVertices adds 2 rays for each polygon vertice
// and sets their direction with a very small offset to the left
//and right of the vertice
func (p *Particle) SetRaysDirToPolyVertices(polygons Polygons) {
	for _, polygon := range polygons {
		for _, vertex := range polygon.Loop {
			rayLeft := NewRay(p.Pos)
			rayLeft.SetDir(vertex.X-0.0001, vertex.Y-0.0001)
			rayRight := NewRay(p.Pos)
			rayRight.SetDir(vertex.X+0.0001, vertex.Y+0.0001)

			p.Rays = append(p.Rays, rayLeft, rayRight)
		}
	}
}

// SortRaysClockwise sorts the rays clockwise by angle
func (p *Particle) SortRaysClockwise() {
	sort.Slice(p.Rays, func(i, j int) bool {
		return p.Rays[i].B.Degrees() < p.Rays[j].B.Degrees()
	})
}

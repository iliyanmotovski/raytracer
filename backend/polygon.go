package backend

import (
	"errors"
	"fmt"
	"math"

	"github.com/iliyanmotovski/raytracer/backend/vector"
)

// Polygon represents a plane which is defined
// by its vertices coordinates
type Polygon struct {
	Loop          vector.Loop
	VerticesCount int
}

// GetBoundaries returns all boundaries (sides) of the polygon
func (p *Polygon) GetBoundaries() Boundaries {
	result := make(Boundaries, p.VerticesCount)

	for i, vertex := range p.Loop {
		next := p.Loop[0]
		if i != len(p.Loop)-1 {
			next = p.Loop[i+1]
		}

		result[i] = &Boundary{vector.Edge{A: vertex, B: next}}
	}

	return result
}

// IsPointContainedInPolygon returns whether a given point is contained
// somewhere inside the Polygon
func (p *Polygon) IsPointContainedInPolygon(point *vector.Vector) bool {
	return p.Loop.IsPointContainedInLoop(point, true)
}

// ContainsVertice returns whether the provided points corresponds
// with some of the Polygon vertices (corners)
func (p *Polygon) ContainsVertice(v *vector.Vector) bool {
	for _, e := range p.Loop {
		if e.X == v.X && e.Y == v.Y {
			return true
		}
	}

	return false
}

// IsConvex checks if the Polygon is convex or not
// Definition: A polygon that has all interior angles less than 180°
func (p *Polygon) IsConvex() bool {
	if len(p.Loop) < 4 {
		return true
	}

	sign := false
	n := len(p.Loop)

	for i, vertice := range p.Loop {
		dx1 := p.Loop[(i+2)%n].X - p.Loop[(i+1)%n].X
		dy1 := p.Loop[(i+2)%n].Y - p.Loop[(i+1)%n].Y
		dx2 := vertice.X - p.Loop[(i+1)%n].X
		dy2 := vertice.Y - p.Loop[(i+1)%n].Y

		zcrossproduct := dx1*dy2 - dy1*dx2

		if i == 0 {
			sign = zcrossproduct > 0
		} else if sign != (zcrossproduct > 0) {
			return false
		}
	}

	return true
}

type Polygons []*Polygon

// Validate checks all polygons and returns whether they intersect
// or some has a point which is outside of the scene, it also checks
// if there is a non-convex polygon
func (ps Polygons) Validate(width, height float64) error {
	vertices := ps.getAllVertices()

	scene := &Polygon{Loop: vector.Loop{
		{X: 0, Y: 0},
		{X: width, Y: 0},
		{X: width, Y: height},
		{X: 0, Y: height},
	}}

	polygons := ps[0:len(ps)]
	polygons = append(polygons, scene)

	for i, polygon := range polygons {
		if !polygon.IsConvex() {
			return errors.New("polygon is not convex")
		}

		for _, vertice := range vertices {
			l := len(polygons) - 1
			contained := polygon.IsPointContainedInPolygon(vertice)

			if i == l && !contained {
				return fmt.Errorf("point X: %v , Y: %v is outside the scene", vertice.X, vertice.Y)
			}

			if i != l && contained && !polygon.ContainsVertice(vertice) {
				return fmt.Errorf("point X: %v , Y: %v is inside another polygon", vertice.X, vertice.Y)
			}
		}
	}

	return nil
}

// getAllVertices returns all vertices of all polygons in an array
func (ps Polygons) getAllVertices() vector.Vectors {
	result := vector.Vectors{}

	for _, polygon := range ps {
		for _, vertice := range polygon.Loop {
			result = append(result, vertice)
		}
	}

	return result
}

// Triangle extends Polygon
type Triangle struct {
	Polygon
}

// Area returns the area of the triangle
func (t *Triangle) Area() float64 {
	a := t.Loop[0]
	b := t.Loop[1]
	c := t.Loop[2]

	return math.Abs((a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y)) / 2)
}

type Triangles []*Triangle

// !!!This function requires the passed vectors to be sorted clockwise by angle!!!
// NewClockwiseTriangleFan takes a center point and an array of vectors which are
// around the center, if we were to visualise the input of this method it would look
// like a circle of dots with a dot in the middle. It connects all the dots with their
// neighbour to the right and with the center dot, the result if we were to visualize
// looks a bit like a "fan" or a "wheel"
func NewClockwiseTriangleFan(center *vector.Vector, edges vector.Loop) Triangles {
	result := Triangles{}

	for i, edge := range edges {
		next := edges[0]
		if i != len(edges)-1 {
			next = edges[i+1]
		}

		result = append(result, &Triangle{Polygon{
			Loop: vector.Loop{
				{X: center.X, Y: center.Y},
				{X: edge.X, Y: edge.Y},
				{X: next.X, Y: next.Y},
			},
			VerticesCount: 3,
		}})
	}

	return result
}

// Area returns the combined area of all triangles
func (ts Triangles) Area() float64 {
	var result float64

	for _, t := range ts {
		result += t.Area()
	}

	return result
}

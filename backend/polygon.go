package backend

import (
	"errors"
	"fmt"
	"math"

	"github.com/iliyanmotovski/raytracer/backend/vector"
)

type Polygon struct {
	Loop          vector.Loop
	VerticesCount int
}

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

func (p *Polygon) IsPointContainedInPolygon(point *vector.Vector) bool {
	return p.Loop.IsPointContainedInLoop(point, true)
}

func (p *Polygon) ContainsVertice(v *vector.Vector) bool {
	for _, e := range p.Loop {
		if e.X == v.X && e.Y == v.Y {
			return true
		}
	}

	return false
}

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

func (ps Polygons) getAllVertices() vector.Vectors {
	result := vector.Vectors{}

	for _, polygon := range ps {
		for _, vertice := range polygon.Loop {
			result = append(result, vertice)
		}
	}

	return result
}

type Triangle struct {
	Polygon
}

func (t *Triangle) Area() float64 {
	a := t.Loop[0]
	b := t.Loop[1]
	c := t.Loop[2]

	return math.Abs((a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y)) / 2)
}

type Triangles []*Triangle

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

func (ts Triangles) Area() float64 {
	var result float64

	for _, t := range ts {
		result += t.Area()
	}

	return result
}

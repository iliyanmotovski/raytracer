package backend

import "github.com/iliyanmotovski/raytracer/backend/vector"

type Ray struct {
	vector.Edge
}

func NewRay(pos *vector.Vector) *Ray {
	return &Ray{Edge: vector.Edge{
		A: pos,
		B: &vector.Vector{X: 1, Y: 0},
	}}
}

func (r *Ray) SetDir(x, y float64) {
	r.B.X = x - r.A.X
	r.B.Y = y - r.A.Y
	*r.B = r.B.Normalize()
}

func (r *Ray) Cast(b *Boundary) (*vector.Vector, bool) {
	return r.Intersection(&b.Edge)
}

type Rays []*Ray

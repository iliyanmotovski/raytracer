package backend

import "github.com/iliyanmotovski/raytracer/backend/vector"

// Ray is a domain level extension of vector.Edge
// it represents a "ray of light" which can be fired
// at any direction
type Ray struct {
	vector.Edge
}

// Creates a new Ray which later can be given direction
func NewRay(pos *vector.Vector) *Ray {
	return &Ray{Edge: vector.Edge{
		// position of the starting point of the ray
		A: pos,
		// direction of the ray
		B: &vector.Vector{X: 1, Y: 0},
	}}
}

// Gives direction to the Ray
func (r *Ray) SetDir(x, y float64) {
	r.B.X = x - r.A.X
	r.B.Y = y - r.A.Y
	*r.B = r.B.Normalize()
}

// Cast casts the ray in the given direction, and checks
// if it intersects with a given boundary, returns
// the point of intersection and boolean if intersection
// occurred
func (r *Ray) Cast(b *Boundary) (*vector.Vector, bool) {
	return r.Intersection(&b.Edge)
}

type Rays []*Ray

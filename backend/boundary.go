package backend

import "github.com/iliyanmotovski/raytracer/backend/vector"

type Boundary struct {
	vector.Edge
}

type Boundaries []*Boundary

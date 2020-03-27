package backend

import "github.com/iliyanmotovski/raytracer/backend/vector"

// Boundary is a domain level wrapper over vector.Edge
// it represents a "solid" line from from point A to B
type Boundary struct {
	vector.Edge
}

type Boundaries []*Boundary

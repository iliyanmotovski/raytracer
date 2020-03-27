package vector

type Edge struct {
	A, B *Vector
}

func (e1 *Edge) Intersection(e2 *Edge) (*Vector, bool) {
	x1 := e2.A.X
	y1 := e2.A.Y
	x2 := e2.B.X
	y2 := e2.B.Y

	x3 := e1.A.X
	y3 := e1.A.Y
	x4 := e1.A.X + e1.B.X
	y4 := e1.A.Y + e1.B.Y

	den := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if den == 0 {
		return nil, false
	}

	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) / den
	u := -((x1-x2)*(y1-y3) - (y1-y2)*(x1-x3)) / den

	if t > 0 && t < 1 && u > 0 {
		x := x1 + t*(x2-x1)
		y := y1 + t*(y2-y1)
		return &Vector{X: x, Y: y}, true
	}

	return nil, false
}

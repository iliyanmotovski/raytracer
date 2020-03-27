package vector

// Loop represents a an array of vectors which form a "loop"
type Loop Vectors

// IsPointContainedInLoop returns whether the provided point is contained inside the loop
// it takes closeLoop argument to know if the loop is closed or not (the last point has the
// same coordinates as the first)
func (l Loop) IsPointContainedInLoop(p *Vector, closeLoop bool) bool {
	if len(l) == 0 {
		return false
	}

	loop := l[0:len(l)]
	if closeLoop {
		loop = append(loop, l[0])
	}

	inside := false
	for i, j := 0, len(loop)-1; i < len(loop); i, j = i+1, i {
		if (loop[i].Y > p.Y) != (loop[j].Y > p.Y) && p.X < (loop[j].X-loop[i].X)*(p.Y-loop[i].Y)/(loop[j].Y-loop[i].Y)+loop[i].X {
			inside = !inside
		}
	}

	return inside
}

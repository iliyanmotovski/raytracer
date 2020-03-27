package vector

import "math"

type Vector struct {
	X, Y float64
}

func (v Vector) Degrees() float64 {
	return math.Atan2(v.Y, v.X) * 180 / math.Pi
}

func (v Vector) Distance(b Vector) float64 { return v.Sub(b).Norm() }

func (v Vector) Sub(b Vector) Vector {
	return Vector{X: v.X - b.X, Y: v.Y - b.Y}
}

func (v Vector) Norm() float64 { return math.Sqrt(v.Dot(v)) }

func (v Vector) Dot(b Vector) float64 {
	return v.X*b.X + v.Y*b.Y
}

func (v Vector) Normalize() Vector {
	return v.MultiplyByScalar(1. / v.Length())
}

func (v Vector) MultiplyByScalar(s float64) Vector {
	return Vector{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

type Vectors []*Vector

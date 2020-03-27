package vector

import "math"

// Vector - struct holding X Y values of a 2D vector
type Vector struct {
	X, Y float64
}

// Degrees returns the vector angle in degrees
func (v Vector) Degrees() float64 {
	return math.Atan2(v.Y, v.X) * 180 / math.Pi
}

// Distance returns the distance from one vector to another
func (v Vector) Distance(b Vector) float64 { return v.Sub(b).Norm() }

// Sub subtracts two vectors
func (v Vector) Sub(b Vector) Vector {
	return Vector{X: v.X - b.X, Y: v.Y - b.Y}
}

// Norm returns the vector norm
func (v Vector) Norm() float64 { return math.Sqrt(v.Dot(v)) }

// Dot returns the dot product of 2 vectors -> cosine of the angle between them
func (v Vector) Dot(b Vector) float64 {
	return v.X*b.X + v.Y*b.Y
}

// Normalize returns the vector normalized
func (v Vector) Normalize() Vector {
	return v.MultiplyByScalar(1. / v.Length())
}

// MultiplyByScalar scales the vector (modifies its length)
func (v Vector) MultiplyByScalar(s float64) Vector {
	return Vector{
		X: v.X * s,
		Y: v.Y * s,
	}
}

// Length returns the length of the vector
func (v Vector) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

type Vectors []*Vector

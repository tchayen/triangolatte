package triangolatte

import "math"

type Point struct {
	X, Y float64
}

// Add adds two two-element vectors.
func (p Point) Add(r Point) Point {
	return Point{p.X + r.X, p.Y + r.Y}
}

// Sub subtracts two two-element vectors.
func (p Point) Sub(r Point) Point {
	return Point{p.X - r.X, p.Y - r.Y}
}

// Scale sets length of vector to given value.
func (p Point) Scale(f float64) Point {
	norm := float64(math.Sqrt(float64(p.X*p.X + p.Y*p.Y)))
	return Point{p.X / norm * f, p.Y / norm * f}
}

// Normalize scales vector to have unit length.
func (p Point) Normalize() Point {
	return p.Scale(1)
}

// Dot calculates dot product of two two-element vectors.
func (p Point) Dot(r Point) float64 {
	return p.X*r.X + p.Y*r.Y
}

// Cross takes one of the approaches to calculating 2D cross product (p.X*r.Y -
// p.Y*r.X).
func (p Point) Cross(r Point) float64 {
	return p.X*r.Y - p.Y*r.X
}

// Distance2 calculates squared distance between two points.
func (p Point) Distance2(r Point) float64 {
	return (p.X-r.X)*(p.X-r.X) + (p.Y-r.Y)*(p.Y-r.Y)
}

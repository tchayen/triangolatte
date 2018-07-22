package triangolatte

import "math"

type Point struct {
	X, Y float64
}

func (p Point) Add(r Point) Point {
	return Point{p.X + r.X, p.Y + r.Y}
}

func (p Point) Sub(r Point) Point {
	return Point{p.X - r.X, p.Y - r.Y}
}

func (p Point) Scale(f float64) Point {
	norm := float64(math.Sqrt(float64(p.X*p.X + p.Y*p.Y)))
	return Point{p.X / norm * f, p.Y / norm * f}
}

func (p Point) Normalize() Point {
	return p.Scale(1)
}

func (p Point) Dot(r Point) float64 {
	return p.X*r.X + p.Y*r.Y
}

func (p Point) Cross(r Point) float64 {
	return p.X*r.Y - p.Y*r.X
}

func (p Point) Distance2(r Point) float64 {
	return (p.X-r.X)*(p.X-r.X) + (p.Y-r.Y)*(p.Y-r.Y)
}

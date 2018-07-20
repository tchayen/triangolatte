package line

import . "triangolatte/pkg/point"

func Normal(points []Point, width int) (triangles []float64) {
	width /= 2.0
	triangles = make([]float64, 0, len(points)*12)
	for i := 0; i <= len(points)-2; i++ {
		dx := points[i+1].X - points[i].X
		dy := points[i+1].Y - points[i].Y
		n1 := Point{dy, -dx}.Scale(float64(width))
		n2 := Point{-dy, dx}.Scale(float64(width))

		v0, v1 := points[i+1].Add(n2).Pair()
		v2, v3 := points[i].Add(n2).Pair()
		v4, v5 := points[i].Add(n1).Pair()
		v6, v7 := points[i].Add(n1).Pair()
		v8, v9 := points[i+1].Add(n1).Pair()
		v10, v11 := points[i+1].Add(n2).Pair()

		triangles = append(triangles, v0, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11)
	}

	return triangles
}

func calculateNormals(x, y float64) [2]Point {
	return [2]Point{
		Point{y, -x}.Normalize(),
		Point{-y, x}.Normalize(),
	}
}

func Miter(points []Point, width int) (triangles []float64) {
	width /= 2.0
	triangles = make([]float64, 0, len(points)*12)
	var dx, dy float64
	var miter1, miter2 Point
	var n1, n2 [2]Point

	dx = points[1].X - points[0].X
	dy = points[1].Y - points[0].Y

	n2 = calculateNormals(dx, dy)
	miter2 = n2[0].Scale(float64(width))

	for i := 1; i < len(points)-1; i++ {
		// Shift calculated values.
		n1 = n2
		miter1 = miter2

		dx = points[i+1].X - points[i].X
		dy = points[i+1].Y - points[i].Y

		n2 = calculateNormals(dx, dy)

		// Find tangent vector to both lines in the middle point.
		tangent := (points[i+1].Sub(points[i])).Normalize().Add((points[i].Sub(points[i-1])).Normalize()).Normalize()

		// Miter vector is perpendicular to the tangent and crosses extensions of
		// normal-translated lines in miter join points.
		unitMiter := Point{-tangent.Y, tangent.X}

		// Length of the miter vector projected onto one of the normals.
		// Choice of normal is arbitrary, each of them would work.
		miterLength := float64(width) / unitMiter.Dot(n1[0])
		miter2 = unitMiter.Scale(miterLength)

		v0, v1 := points[i].Sub(miter2).Pair()
		v2, v3 := points[i-1].Sub(miter1).Pair()
		v4, v5 := points[i-1].Add(miter1).Pair()
		v6, v7 := points[i-1].Add(miter1).Pair()
		v8, v9 := points[i].Add(miter2).Pair()
		v10, v11 := points[i].Sub(miter2).Pair()

		triangles = append(triangles, v0, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11)
	}

	// Use last normal as another 'neutral element' for miter join.
	n := len(points)
	lastMiter := n2[0].Scale(float64(width))

	v0, v1 := points[n-1].Sub(lastMiter).Pair()
	v2, v3 := points[n-2].Sub(miter1).Pair()
	v4, v5 := points[n-2].Add(miter1).Pair()
	v6, v7 := points[n-2].Add(miter1).Pair()
	v8, v9 := points[n-1].Add(lastMiter).Pair()
	v10, v11 := points[n-1].Sub(lastMiter).Pair()

	triangles = append(triangles, v0, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11)

	return triangles
}

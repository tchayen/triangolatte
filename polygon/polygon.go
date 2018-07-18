package polygon

import (
	"sort"
	. "triangolate/point"
)

type Set map[int]bool

type Triangle struct {
	A, B, C Point
}

func Cyclic(i, n int) int {
	return (i%n + n) % n
}

func IsReflex(a, b, c Point) bool {
	return b.Sub(a).Cross(c.Sub(b)) < 0.0
}

func SameSide(p1, p2, a, b Point) bool {
	return b.Sub(a).Cross(p1.Sub(a))*b.Sub(a).Cross(p2.Sub(a)) >= 0
}

func IsInsideTriangle(t Triangle, p Point) bool {
	return SameSide(p, t.A, t.B, t.C) &&
		SameSide(p, t.B, t.A, t.C) &&
		SameSide(p, t.C, t.A, t.B)
}

func SplitConvexAndReflex(points []Point) (convex, reflex Set) {
	n := len(points)
	convex, reflex = make(Set, len(points)), make(Set, len(points))

	for i := 0; i < n; i++ {
		if IsReflex(points[Cyclic(i-1, n)], points[i], points[Cyclic(i+1, n)]) {
			reflex[i] = true
		} else {
			convex[i] = true
		}
	}
	return convex, reflex
}

func DetectEars(points []Point, reflex Set) (ears Set) {
	n := len(points)
	ears = make(Set, len(points))

	for i := 0; i < n; i++ {
		if reflex[i] {
			continue
		}

		isEar := true
		for j := 0; j < n; j++ {
			// It is ok to skip reflex vertices and the ones that actually belong to
			// the triangle.
			if reflex[j] || j == Cyclic(i-1, n) || j == i || j == Cyclic(i+1, n) {
				continue
			}

			// If triangle contains points[j], points[i] cannot be an ear tip.
			if IsInsideTriangle(Triangle{
				points[Cyclic(i-1, n)],
				points[i],
				points[Cyclic(i+1, n)],
			}, points[j]) {
				isEar = false
			}
		}
		if isEar {
			ears[i] = true
		}
	}
	return ears
}

func CombinePolygons(outer, inner []Point) (result []Point) {
	xMax := 0.0
	mIndex := 0
	for i := 0; i < len(inner); i++ {
		if inner[i].X > xMax {
			xMax = inner[i].X
			mIndex = i
		}
	}

	m := inner[mIndex]

	// Find the edges that intersect with ray `M + t * (1, 0)`. Let `K` be the
	// closest visible point to `M` on this ray.
	var k Point
	foundK := false
	var k1, k2 int
	for i, j := len(outer)-1, 0; j < len(outer); i, j = j, j+1 {
		// Skip edges that does not have their first point below `M` and the second
		// one above.
		if outer[i].Y > m.Y || outer[j].Y < m.Y {
			continue
		}

		// Calculate simplified intersection of ray (1, 0) and [V_i, V_j] segment.
		v1 := m.Sub(outer[i])
		v2 := outer[j].Sub(outer[i])

		t1 := v2.Cross(v1) / v2.Y
		t2 := v1.Y / v2.Y

		if t1 >= 0.0 && t2 >= 0.0 && t2 <= 1.0 {
			// If there is no current `k` candidate or this one is closer.
			if !foundK || t1-m.X < k.X {
				k = Point{X: t1 + m.X, Y: m.Y}
				k1, k2 = i, j
				foundK = true
			}
		} else {
			panic("Cannot calculate intersection, problematic data")
		}
	}

	var visibleIndex, pIndex int = -1, -1

	// If `K` is vertex of the outer polygon, `M` and `K` are mutually visible.
	for i := 0; i < len(outer); i++ {
		if outer[i] == k {
			visibleIndex = i
		}
	}

	// Otherwise, `K` is an interior point of the edge `[V_k_1, V_k_2]`. Find `P`
	// which is endpoint with greater x-value.
	if outer[k1].X > outer[k2].X {
		pIndex = k1
	} else {
		pIndex = k2
	}

	// Check with all vertices of the outer polygon to be outside of the
	// triangle `[M, K, P]`. If it is true, `M` and `P` are mutually visible.
	anyInside := false
	for i := 0; i < len(outer); i++ {
		anyInside = anyInside && (outer[i] == outer[pIndex] || !IsInsideTriangle(Triangle{m, k, outer[pIndex]}, outer[i]))
	}

	if visibleIndex < 0 && anyInside {
		visibleIndex = pIndex
	}

	// Otherwise at least one reflex vertex lies in `[M, K, P]`. Search for the
	// array of reflex vertices `R` that minimizes the angle between `(1, 0)` and
	// line segment `[M, R]`. If there is exactly one vertex in `R` then they are
	// mutually visible. If there are multiple such vertices, pick the one closest
	// to `M`.
	if visibleIndex < 0 {
		var reflex []int
		n := len(outer)
		for i := 0; i < n; i++ {
			if !IsInsideTriangle(Triangle{m, k, outer[pIndex]}, outer[i]) || !IsReflex(outer[Cyclic(i-1, n)], outer[i], outer[Cyclic(i+1, n)]) {
				continue
			}
			reflex = append(reflex, i)
		}
		sort.Ints(reflex)
		visibleIndex = reflex[0]
	}

	if visibleIndex < 0 {
		panic("Could not find visible vertex")
	}

	result = make([]Point, 0, len(outer)+len(inner)+2)
	result = append(outer[:visibleIndex])
	result = append(inner[:])
	result = append(result, inner[mIndex], outer[visibleIndex])
	result = append(outer[visibleIndex:])

	return result
}

type byMaxX [][]Point

func (polygons byMaxX) Len() int {
	return len(polygons)
}

func (polygons byMaxX) Swap(i, j int) {
	polygons[i], polygons[j] = polygons[j], polygons[i]
}

func (polygons byMaxX) Less(i, j int) bool {
	xMax := 0.0

	for k := 0; k < len(polygons[i]); k++ {
		if polygons[i][k].X > xMax {
			xMax = polygons[i][k].X
		}
	}

	for k := 0; k < len(polygons[j]); k++ {
		if polygons[j][k].X > xMax {
			return true
		}
	}
	return false
}

func EliminateHoles(outer []Point, inners [][]Point) {
	sort.Sort(byMaxX(inners))

	var inner []Point
	for len(inners) > 0 {
		inner, inners = inners[0], inners[1:]
		outer = CombinePolygons(outer, inner)
	}
}

func EarCut(points []Point) (triangles []float64) {
	n := len(points)
	if n < 3 {
		panic("Cannot triangulate less than three points")
	}

	var _, reflex Set = SplitConvexAndReflex(points)
	var ears Set = DetectEars(points, reflex)

	triangles = make([]float64, 0, (n-1)/2)
	for len(points) > 3 {
		i := 0
		for k := range ears {
			i = k
		}

		v1, v2 := points[Cyclic(i-1, n)].Pair()
		v3, v4 := points[i].Pair()
		v5, v6 := points[Cyclic(i+1, n)].Pair()

		triangles = append(triangles, v1, v2, v3, v4, v5, v6)
		points = append(points[:i], points[i+1:]...)
		n = n - 1

		_, reflex = SplitConvexAndReflex(points)
		ears = DetectEars(points, reflex)
	}
	return triangles
}

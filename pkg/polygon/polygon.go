package polygon

import (
	"sort"
	. "triangolatte/pkg/point"
	"errors"
	"container/list"
)

type Set map[int]bool

type Triangle struct {
	A, B, C Point
}

func cyclic(i, n int) int {
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

func splitConvexAndReflex(points []Point, indexMap []int) (convex, reflex Set) {
	n := len(indexMap)
	convex, reflex = make(Set, len(points)), make(Set, len(points))

	for i := 0; i < n; i++ {
		a := points[indexMap[cyclic(i-1, n)]]
		b := points[indexMap[i]]
		c := points[indexMap[cyclic(i+1, n)]]

		if IsReflex(a, b, c) {
			reflex[i] = true
		} else {
			convex[i] = true
		}
	}
	return convex, reflex
}

func detectEars(points []Point, reflex Set, indexMap []int) *list.List {
	n := len(indexMap)
	ears := list.New()

	for i := 0; i < n; i++ {
		if reflex[i] {
			continue
		}

		isEar := true
		for j := 0; j < n; j++ {
			// It is ok to skip reflex vertices and the ones that actually belong to
			// the triangle.
			if !reflex[j] || j == cyclic(i-1, n) || j == i || j == cyclic(i+1, n) {
				continue
			}

			// If triangle contains points[j], points[i] cannot be an ear tip.
			if IsInsideTriangle(Triangle{
				points[indexMap[cyclic(i-1, n)]],
				points[indexMap[i]],
				points[indexMap[cyclic(i+1, n)]],
			}, points[indexMap[j]]) {
				isEar = false
			}
		}

		if isEar {
			ears.PushBack(i)
		}
	}
	return ears
}

func combinePolygons(outer, inner []Point) ([]Point, error) {
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
			return nil, errors.New("cannot calculate intersection, problematic data")
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
	allOutside := true
	for i := range outer {
		// We have to skip M, K and P vertices. Since M is from the inner
		// polygon and K was proved to not match any vertex, the only one to
		// check is pIndex
		if i == pIndex {
			continue
		}

		if IsInsideTriangle(Triangle{m, k, outer[pIndex]}, outer[i]) {
			allOutside = false
		}
	}

	if visibleIndex < 0 && allOutside {
		visibleIndex = pIndex
	}

	// Otherwise at least one reflex vertex lies in `[M, K, P]`. Search for the
	// array of reflex vertices `R` that minimizes the angle between `(1, 0)` and
	// line segment `[M, R]`. If there is exactly one vertex in `R` then they are
	// mutually visible. If there are multiple such vertices, pick the one closest
	// to `M`.
	if visibleIndex < 0 {
		reflex := list.New()
		n := len(outer)
		for i := 0; i < n; i++ {
			notInside := !IsInsideTriangle(Triangle{m, k, outer[pIndex]}, outer[i])
			notReflex := !IsReflex(outer[cyclic(i-1, n)], outer[i], outer[cyclic(i+1, n)])
			if notInside || notReflex {
				continue
			}
			reflex.PushBack(i)
		}
		var closest int
		var maxDist float64

		for r := reflex.Front(); r != nil; r = r.Next() {
			i := r.Value.(int)
			dist := outer[i].Distance2(outer[closest])
			if dist > maxDist {
				closest = i
				maxDist = dist
			}
		}
		visibleIndex = closest
	}

	if visibleIndex < 0 {
		return nil, errors.New("could not find visible vertex")
	}

	n := len(inner)
	result := make([]Point, 0, len(outer)+len(inner)+2)
	result = append(result, outer[:visibleIndex+1]...)
	for i := 0; i < n; i++ {
		result = append(result, inner[cyclic(mIndex+i, n)])
	}
	result = append(result, inner[mIndex], outer[visibleIndex])
	result = append(result, outer[visibleIndex+1:]...)

	return result, nil
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
			return false
		}
	}
	return true
}

func eliminateHoles(outer []Point, inners [][]Point) ([]Point, error) {
	sort.Sort(byMaxX(inners))

	var inner []Point
	for len(inners) > 0 {
		inner, inners = inners[0], inners[1:]
		var err error
		outer, err = combinePolygons(outer, inner)
		if err != nil {
			return nil, err
		}
	}
	return outer, nil
}

func EarCut(points []Point, holes [][]Point) ([]float64, error) {
	n := len(points)
	if n < 3 {
		return nil, errors.New("cannot triangulate less than three points")
	}

	if len(holes) > 0 {
		var err error
		points, err = eliminateHoles(points, holes)
		if err != nil {
			return nil, err
		}
	}

	var indexMap = make([]int, n)
	for i := range indexMap {
		indexMap[i] = i
	}

	var _, reflex Set = splitConvexAndReflex(points, indexMap)
	var ears = detectEars(points, reflex, indexMap)

	// Any triangulation of simple polygon has `n-2` triangles.
	t, triangles := 0, make([]float64, (n-2) * 6)
	for len(indexMap) > 3 {

		if ears.Len() == 0 {
			return nil, errors.New("could not detect any ear tip in a non-empty polygon")
		}
		i := ears.Remove(ears.Front()).(int)

		triangles[t], triangles[t+1] = points[indexMap[cyclic(i-1, n)]].Pair()
		triangles[t+2], triangles[t+3] = points[indexMap[i]].Pair()
		triangles[t+4], triangles[t+5] = points[indexMap[cyclic(i+1, n)]].Pair()
		t += 6

		indexMap = append(indexMap[:i], indexMap[i+1:]...) // Skip `points[indexMap[i]]`.
		n = n - 1

		_, reflex = splitConvexAndReflex(points, indexMap)
		ears = detectEars(points, reflex, indexMap)
	}

	triangles[t], triangles[t+1] = points[indexMap[0]].Pair()
	triangles[t+2], triangles[t+3] = points[indexMap[1]].Pair()
	triangles[t+4], triangles[t+5] = points[indexMap[2]].Pair()

	return triangles, nil
}

package polygon

import (
	"sort"
	. "triangolatte/pkg/point"
	"triangolatte/pkg/cyclicList"
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

func setReflex(points *cyclicList.CyclicList) {
	n := points.Len()

	for i, p := 0, points.Front(); i < n; i, p = i + 1, p.Next() {
		if IsReflex(p.Prev().Point, p.Point, p.Next().Point) {
			p.Reflex = true
		}
	}
}

func isEar(p *cyclicList.Element, t Triangle) bool {
	n := p.List.Len()
	for i, r := 0, p.List.Front(); i < n; i, r = i + 1, r.Next() {
		// It is ok to skip reflex vertices and the ones that actually belong to
		// the triangle.
		if p.Reflex || r == p.Prev() || r == p || r == p.Next() {
			continue
		}

		// If triangle contains points[j], points[i] cannot be an ear tip.
		if IsInsideTriangle(t, r.Point) {
			return false
		}
	}
	return true
}

func detectEars(points *cyclicList.CyclicList) *list.List {
	ears := list.New()

	n := points.Len()
	for i, p := 0, points.Front(); i < n; i, p = i + 1, p.Next() {
		if p.Reflex {
			continue
		}

		t := Triangle{p.Prev().Point, p.Point, p.Next().Point}
		if isEar(p, t) {
			ears.PushBack(p)
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

	c := cyclicList.NewFromArray(points)

	if len(holes) > 0 {
		var err error
		points, err = eliminateHoles(points, holes)
		if err != nil {
			return nil, err
		}
	}

	setReflex(c)
	var ears = detectEars(c)

	// Any triangulation of simple polygon has `n-2` triangles.
	i, t := 0, make([]float64, (n-2) * 6)
	for c.Len() > 3 {
		if ears.Len() == 0 {
			return nil, errors.New("could not detect any ear tip in a non-empty polygon")
		}

		ear := ears.Remove(ears.Front()).(*cyclicList.Element)

		t[i+0], t[i+1] = ear.Prev().Point.Pair()
		t[i+2], t[i+3] = ear.Point.Pair()
		t[i+4], t[i+5] = ear.Next().Point.Pair()
		i += 6

		// Skip `points[indexMap[i]]`.
		c.Remove(ear)
		n = n - 1

		// // If an adjacent vertex is convex, it remains convex.
		// // If an adjacent vertex is an ear, it does not necessarily remain an ear.
		// // If the adjacent vertex is reflex, it is possible that it becomes
		// // convex and, possibly, an ear.
		// check := func(index int) {
		// 	triangle := Triangle{
		// 		points[indexMap[cyclic(index-1, n)]],
		// 		points[indexMap[index]],
		// 		points[indexMap[cyclic(index+1, n)]],
		// 	}
        //
		// 	if reflex[index] {
		// 		if !IsReflex(triangle.A, triangle.B, triangle.C) {
		// 			delete(reflex, index)
        //
		// 			if isEar(index, triangle, points, reflex, indexMap) {
		// 				for e := ears.Front(); e != nil; e = e.Next() {
		// 					if e.Value.(int) < index && e.Next().Value.(int) > index {
		// 						ears.InsertAfter(index, e)
		// 					}
		// 				}
		// 			}
		// 		}
		// 	} else {
		// 		for e := ears.Front(); e != nil; e = e.Next() {
		// 			if e.Value.(int) == index {
		// 				if !isEar(index, triangle, points, reflex, indexMap) {
		// 					ears.Remove(e)
		// 				}
		// 			}
		// 		}
		// 	}
		// }
		// check(ear - 1)
		// check(ear)

		setReflex(c)
		ears = detectEars(c)
	}

	p := c.Front()
	t[i+0], t[i+1] = p.Point.Pair(); p = p.Next()
	t[i+2], t[i+3] = p.Point.Pair(); p = p.Next()
	t[i+4], t[i+5] = p.Point.Pair(); p = p.Next()

	return t, nil
}

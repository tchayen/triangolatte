package triangolatte

import (
	"container/list"
	"errors"
	"sort"
)

// Set of int values.
type Set map[int]bool

func cyclic(i, n int) int {
	return (i%n + n) % n
}

func triangleArea(a, b, c Point) float64 {
	return a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y)
}

// isReflex checks if angle created by points a, b and c is reflex.
//
// Angle equal to math.Pi is considered convex for practical reasons (it can be
// used just fine in the triangulation).
//
// But generally, math.Pi angle should not happen since collinear points are
// redundant and therefore they should be eliminated in preprocessing.
func isReflex(a, b, c Point) bool {
	return (b.X-a.X)*(c.Y-b.Y)-(c.X-b.X)*(b.Y-a.Y) < 0
}

// isInsideTriangle checks if given point P lays inside triangle [A, B, C].
func isInsideTriangle(a, b, c, p Point) bool {
	return (c.X-p.X)*(a.Y-p.Y)-(a.X-p.X)*(c.Y-p.Y) >= 0 &&
		(a.X-p.X)*(b.Y-p.Y)-(b.X-p.X)*(a.Y-p.Y) >= 0 &&
		(b.X-p.X)*(c.Y-p.Y)-(c.X-p.X)*(b.Y-p.Y) >= 0
}

func setReflex(points *CyclicList) {
	n := points.Len()

	for i, p := 0, points.Front(); i < n; i, p = i+1, p.Next() {
		if isReflex(p.Prev().Point, p.Point, p.Next().Point) {
			p.Reflex = true
		} else {
			p.Reflex = false
		}
	}
}

func isEar(p *Element) bool {
	// Skip given point and its neighbours.
	r := p.Next().Next()
	for r != p.Prev() {
		if p.Reflex {
			continue
		}

		// If triangle contains points[j], points[i] cannot be an ear tip.
		a, b, c := p.Prev().Point, p.Point, p.Next().Point
		if isInsideTriangle(a, b, c, r.Point) {
			return false
		}
		r = r.Next()
	}
	return true
}

func detectEars(points *CyclicList) *list.List {
	ears := list.New()

	n := points.Len()
	for i, p := 0, points.Front(); i < n; i, p = i+1, p.Next() {
		if p.Reflex {
			continue
		}

		if isEar(p) {
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

		if isInsideTriangle(m, k, outer[pIndex], outer[i]) {
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
			notInside := !isInsideTriangle(m, k, outer[pIndex], outer[i])
			notReflex := !isReflex(outer[cyclic(i-1, n)], outer[i], outer[cyclic(i+1, n)])
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

// If an adjacent vertex is convex, it remains convex.
// If an adjacent vertex is an ear, it does not necessarily remain an ear.
// If an adjacent vertex is reflex, it is possible that it becomes
// convex and, possibly, an ear.
func checkVertex(element *Element, ears *list.List) {
	a, b, c := element.Prev().Point, element.Point, element.Next().Point

	if element.Reflex {
		if !isReflex(a, b, c) {
			element.Reflex = false

			if isEar(element) {
				element.Ear = ears.PushBack(element)
			}
		}
	} else {
		if element.Ear != nil && !isEar(element) {
			ears.Remove(element.Ear)
			element.Ear = nil
		}
	}
}

// EarCut triangulates given CCW polygon using ear clipping algorithm (takes
// O(n^2) time).
//
// Produces  array of two-coordinate, CCW triangles, put one after another.
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

	remainingPoints := NewFromArray(points)

	setReflex(remainingPoints)
	var ears = detectEars(remainingPoints)

	// Any triangulation of simple polygon has `n-2` triangles.
	i, t := 0, make([]float64, (n-2)*6)
	for remainingPoints.Len() > 3 {
		if ears.Len() == 0 {
			return nil, errors.New("could not detect any more ear tips")
		}

		ear := ears.Remove(ears.Front()).(*Element)
		ear.Ear = nil

		a, b, c := ear.Prev().Point, ear.Point, ear.Next().Point

		// Skip adding triangles in which all points are collinear.
		if triangleArea(a, b, c) > 0 {
			t[i+0], t[i+1] = a.X, a.Y
			t[i+2], t[i+3] = b.X, b.Y
			t[i+4], t[i+5] = c.X, c.Y
			i += 6
		}

		// Save previous and next before Remove() gets rid of pointers.
		prev := ear.Prev()
		next := ear.Next()

		remainingPoints.Remove(ear)
		n = n - 1

		checkVertex(prev, ears)
		checkVertex(next, ears)
	}

	p := remainingPoints.Front()
	t[i+0], t[i+1] = p.Point.X, p.Point.Y
	p = p.Next()
	t[i+2], t[i+3] = p.Point.X, p.Point.Y
	p = p.Next()
	t[i+4], t[i+5] = p.Point.X, p.Point.Y

	return t, nil
}

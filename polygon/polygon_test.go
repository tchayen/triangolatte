package polygon

import (
	"math"
	"sort"
	"testing"
	. "triangolatte/point"
)

var vertices = []Point{{50, 110}, {150, 30}, {240, 115}, {320, 65}, {395, 170}, {305, 160}, {265, 240}, {190, 100}, {95, 125}, {100, 215}}

func checkArray(t *testing.T, result, expected []int) {
	if len(result) != len(expected) {
		t.Error("Array sizes don't match")
	}

	for i := 0; i < len(result); i++ {
		if math.Abs(float64(result[i]-expected[i])) > 0.001 {
			t.Error("Value error beyond floating point precision")
		}
	}
}

// TODO: generalise function above
func checkPointArray(t *testing.T, result, expected []Point) {
	if len(result) != len(expected) {
		t.Error("Array sizes don't match")
	}

	for i := 0; i < len(result); i++ {
		if math.Abs(result[i].X-expected[i].X) > 0.001 && math.Abs(result[i].Y-expected[i].Y) > 0.001 {
			t.Error("Value error beyond floating point precision")
		}
	}
}

func TestCyclic(t *testing.T) {
	if cyclic(1, 5) != 1 || cyclic(4, 5) != 4 || cyclic(6, 5) != 1 || cyclic(-1, 5) != 4 || cyclic(-5, 5) != 0 || cyclic(-6, 5) != 4 {
		t.Error("cyclic is broken")
	}
}

func TestIsReflex(t *testing.T) {
	if IsReflex(Point{0, 0}, Point{1, 1}, Point{2, 0}) != true || IsReflex(Point{0, 0}, Point{1, 0}, Point{1, 1}) != false {
		t.Error("IsReflex is broken")
	}
}

func TestSameSide(t *testing.T) {
	if SameSide(Point{3, 1}, Point{4, 2}, Point{0, 0}, Point{5, 3}) != true {
		t.Error("SameSide is broken")
	}
}

func TestIsInsideTriangle(t *testing.T) {
	if IsInsideTriangle(Triangle{vertices[0], vertices[1], vertices[2]}, vertices[7]) != true || IsInsideTriangle(Triangle{vertices[0], vertices[1], vertices[2]}, vertices[5]) != false {
		t.Error("IsInsideTriangle is broken")
	}
}

func TestSplitConvexAndReflex(t *testing.T) {
	convex, reflex := splitConvexAndReflex([]Point{{0, 0}, {2, 3}, {4, 2}, {0, 7}})
	t.Log(convex)
	t.Log(reflex)

	if !(convex[0] && convex[2] && convex[3] && !convex[1] && reflex[1] && !reflex[2]) {
		t.Error("splitConvexAndReflex is broken")
	}
}

func TestDetectEars(t *testing.T) {
	_, reflex := splitConvexAndReflex(vertices)
	earsMap := detectEars(vertices, reflex)
	ears := make([]int, 0, len(earsMap))
	for _, v := range ears {
		ears = append(ears, v)
	}
	sort.Ints(ears)

	expectedEars := []int{3, 4, 6, 9}

	t.Log(earsMap)
	t.Log(expectedEars)
	checkArray(t, ears, expectedEars)

}

func TestEliminateHoles(t *testing.T) {
	polygon := []Point{{0, 30}, {20, 0}, {80, 0}, {90, 40}, {30, 70}}
	holes := [][]Point{
		{{20, 10}, {20, 40}, {50, 40}},
		{{60, 30}, {70, 20}, {50, 10}},
	}

	polygonWithEliminatedHoles := []Point{
		{0, 30}, {20, 0}, {80, 0},
		{90, 40}, {70, 20}, {50, 10},
		{60, 30}, {70, 20}, {90, 40},
		{50, 40}, {20, 10}, {20, 40},
		{50, 40}, {90, 40}, {30, 70},
	}

	eliminated := eliminateHoles(polygon, holes)
	checkPointArray(t, eliminated, polygonWithEliminatedHoles)
}

func TestEarCut(t *testing.T) {

}

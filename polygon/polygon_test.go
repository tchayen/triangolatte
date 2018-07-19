package polygon

import (
	"math"
	"sort"
	"testing"
	. "triangolatte/point"
)

var vertices = []Point{{50, 110}, {150, 30}, {240, 115}, {320, 65}, {395, 170}, {305, 160}, {265, 240}, {190, 100}, {95, 125}, {100, 215}}

func checkIntArray(t *testing.T, result, expected []int) {
	if len(result) != len(expected) {
		t.Error("Array sizes don't match")
	}

	for i := 0; i < len(result); i++ {
		if math.Abs(float64(result[i]-expected[i])) > 0.001 {
			t.Error("Value error beyond floating point precision problem")
		}
	}
}

func checkFloat64Array(t *testing.T, result, expected []float64) {
	if len(result) != len(expected) {
		t.Error("Array sizes don't match")
	}

	for i := 0; i < len(result); i++ {
		if math.Abs(result[i]-expected[i]) > 0.001 {
			t.Error("Value error beyond floating point precision problem")
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
			t.Error("Value error beyond floating point precision problem")
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
	indexMap := []int{0, 1, 2, 3}

	convex, reflex := splitConvexAndReflex([]Point{{0, 0}, {2, 3}, {4, 2}, {0, 7}}, indexMap)
	t.Log(convex)
	t.Log(reflex)

	if !(convex[0] && convex[2] && convex[3] && !convex[1] && reflex[1] && !reflex[2]) {
		t.Error("splitConvexAndReflex is broken")
	}
}

func TestDetectEars(t *testing.T) {
	var indexMap []int = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	_, reflex := splitConvexAndReflex(vertices, indexMap)
	earsMap := detectEars(vertices, reflex, indexMap)
	ears := make([]int, 0, len(earsMap))
	for k, _ := range earsMap {
		ears = append(ears, k)
	}
	sort.Ints(ears)

	expectedEars := []int{3, 4, 6, 9}

	t.Log(ears)
	t.Log(expectedEars)
	checkIntArray(t, ears, expectedEars)
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

func TestEliminateHolesWithNoDirectlyVisible(t *testing.T) {
	// TODO
}

func TestEliminateHolesWithNoPossibleVisibleVertex(t *testing.T) {
	// TODO
}

func TestEarCut(t *testing.T) {
	result := EarCut(vertices, [][]Point{})
	expected := []float64{240, 115, 320, 65, 395, 170, 240, 115, 395, 170, 305, 160, 240, 115, 305, 160, 265, 240, 240, 115, 265, 240, 190, 100, 150, 30, 240, 115, 190, 100, 50, 110, 150, 30, 190, 100, 50, 110, 190, 100, 95, 125, 50, 110, 95, 125, 100, 215}

	t.Log(result)
	t.Log(expected)

	checkFloat64Array(t, result, expected)
}

func TestIncorrectEarCut(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	EarCut([]Point{{0, 0}}, [][]Point{})
}

func TestSortingByXMax(t *testing.T) {
	inners := [][]Point{
		[]Point{{1, 2}},
		[]Point{{0, 0}},
	}
	sort.Sort(byMaxX(inners))
}

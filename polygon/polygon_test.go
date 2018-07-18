package polygon

import (
	"testing"
	. "triangolate/point"
)

var vertices = []Point{{50, 110}, {150, 30}, {240, 115}, {320, 65}, Point{395, 170}, {305, 160}, {265, 240}, {190, 100}, {95, 125}, {100, 215}}

func TestCyclic(t *testing.T) {
	if Cyclic(1, 5) != 1 || Cyclic(4, 5) != 4 || Cyclic(6, 5) != 1 || Cyclic(-1, 5) != 4 || Cyclic(-5, 5) != 0 || Cyclic(-6, 5) != 4 {
		t.Error("Cyclic is broken")
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
	convex, reflex := SplitConvexAndReflex([]Point{{0, 0}, {2, 3}, {4, 2}, {0, 7}})
	t.Log(convex)
	t.Log(reflex)

	if !(convex[0] && convex[2] && convex[3] && !convex[1] && reflex[1] && !reflex[2]) {
		t.Error("SplitConvexAndReflex is broken")
	}
}

func TestDetectEars(t *testing.T) {
	// TODO
}

func TestEliminateHoles(t *testing.T) {
	var polygon = []Point{{0, 30}, {20, 0}, {80, 0}, {90, 40}, {30, 70}}
	var holes = [][]Point{
		{{20, 10}, {20, 40}, {50, 40}},
		{{60, 30}, {70, 20}, {50, 10}},
	}

	var polygonWithEliminatedHoles = []Point{
		{0, 30}, {20,  0}, {80,  0},
		{90, 40}, {70, 20}, {50, 10},
		{60, 30}, {70, 20}, {90, 40},
		{50, 40}, {20, 10}, {20, 40},
		{50, 40}, {90, 40}, {30, 70},
	}

	var withoutHoles []Point
	copy(polygonWithEliminatedHoles, withoutHoles)

	EliminateHoles(polygon, holes)
	for i := 0; i < len(polygonWithEliminatedHoles); i++ {
		if polygonWithEliminatedHoles[i] != withoutHoles[i] {
			t.Error("Incorrect hole elimination")
		}
	}
}

func TestEarCut(t *testing.T) {

}
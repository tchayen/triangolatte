package triangolatte

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

var vertices = []Point{{50, 110}, {150, 30}, {240, 115}, {320, 65}, {395, 170}, {305, 160}, {265, 240}, {190, 100}, {95, 125}, {100, 215}}

func TestPolygonArea(t *testing.T) {
	points := []Point{{2, 2}, {11, 2}, {9, 7}, {4, 10}}
	area := polygonArea(points)

	if area != 45.5 {
		t.Error("polygonArea implementation is wrong")
	}
}

func TestDeviation(t *testing.T) {
	data := []Point{{0, 4}, {3, 1}, {8, 2}, {9, 5}, {4, 6}}
	triangles := []float64{4, 6, 0, 4, 3, 1, 4, 6, 3, 1, 8, 2, 8, 2, 9, 5, 4, 6}

	real, actual, deviation := deviation(data, triangles)
	if deviation > 0 {
		t.Errorf("real: %f, actual: %f", real, actual)
	}
}

func checkFloat64Array(t *testing.T, result, expected []float64) {
	if len(result) != len(expected) {
		t.Error("Array sizes don't match")
	}

	for i, r := range result {
		if math.Abs(r-expected[i]) > 0.001 {
			t.Error("Value error beyond floating point precision problem")
		}
	}
}

func checkBoolArray(t *testing.T, result, expected []bool) {
	if len(result) != len(expected) {
		t.Error("Array sizes don't match")
	}

	for i, r := range result {
		if r != expected[i] {
			t.Error("Value error")
		}
	}
}

func checkPointArray(t *testing.T, result, expected []Point) {
	if len(result) != len(expected) {
		t.Error("Array sizes don't match")
	}

	for i, r := range result {
		if math.Abs(r.X-expected[i].X) > 0.001 && math.Abs(result[i].Y-expected[i].Y) > 0.001 {
			t.Error("Value error beyond floating point precision problem")
		}
	}
}

func TestCyclic(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		r := cyclic(1, 5)
		if r != 1 {
			t.Errorf("%d != %d", r, 1)
		}
	})

	t.Run("overflow", func(t *testing.T) {
		r := cyclic(6, 5)
		if r != 1 {
			t.Errorf("%d != %d", r, 1)
		}
	})

	t.Run("negative overflow", func(t *testing.T) {
		r := cyclic(-1, 5)
		if r != 4 {
			t.Errorf("%d != %d", r, 4)
		}
	})

	t.Run("full loop", func(t *testing.T) {
		r := cyclic(-5, 5)
		if r != 0 {
			t.Errorf("%d != %d", r, 0)
		}
	})

	t.Run("double negative overflow", func(t *testing.T) {
		r := cyclic(-6, 5)
		if r != 4 {
			t.Errorf("%d != %d", r, 4)
		}
	})
}

func TestIsReflex(t *testing.T) {
	t.Run("convex", func(t *testing.T) {
		convex := []Point{{0, 1}, {1, 0}, {2, 1}}
		if isReflex(convex[0], convex[1], convex[2]) {
			t.Error("isReflex: false negative")
		}
	})

	t.Run("reflex", func(t *testing.T) {
		reflex := []Point{{0, 0}, {0, 3}, {2, 3}}
		if !isReflex(reflex[0], reflex[1], reflex[2]) {
			t.Error("isReflex: false positive")
		}
	})

	t.Run("square", func(t *testing.T) {
		square := []Point{{1, 1}, {0, 1}, {0, 0}}
		if isReflex(square[0], square[1], square[2]) {
			t.Error("isReflex: false negative")
		}
	})

	t.Run("another reflex", func(t *testing.T) {
		anotherReflex := []Point{{0, 0}, {2, 3}, {4, 2}}
		if !isReflex(anotherReflex[0], anotherReflex[1], anotherReflex[2]) {
			t.Error("isReflex: false positive")
		}
	})
}

func BenchmarkIsReflex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isReflex(Point{0, 0}, Point{1, 1}, Point{2, 0})
	}
}

func TestIsInsideTriangle(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		if isInsideTriangle(vertices[0], vertices[8], vertices[9], vertices[7]) {
			t.Error("isInsideTriangle is broken")
		}
	})

	t.Run("case 2", func(t *testing.T) {
		if !isInsideTriangle(vertices[0], vertices[1], vertices[5], vertices[7]) {
			t.Error("isInsideTriangle is broken")
		}
	})
}

func BenchmarkIsInsideTriangle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isInsideTriangle(vertices[0], vertices[1], vertices[2], vertices[3])
	}
}

func TestSetReflex(t *testing.T) {
	points := []Point{{0, 0}, {2, 3}, {4, 2}, {0, 7}}
	c := NewFromArray(points)
	setReflex(c)

	result := make([]bool, 4)
	for i, p := 0, c.Front(); i < c.Len(); i, p = i+1, p.Next() {
		result[i] = p.Reflex
	}
	expected := []bool{false, true, false, false}

	t.Log(result, expected)
	checkBoolArray(t, result, expected)
}

func TestDetectEars(t *testing.T) {
	c := NewFromArray(vertices)
	setReflex(c)
	earList := detectEars(c)

	ears := make([]Point, earList.Len())
	for i, e := 0, earList.Front(); e != nil; i, e = i+1, e.Next() {
		ears[i] = e.Value.(*Element).Point
	}

	expectedEars := []Point{vertices[3], vertices[4], vertices[6], vertices[9]}
	checkPointArray(t, ears, expectedEars)
}

func BenchmarkDetectEars(b *testing.B) {
	b.StopTimer()
	c := NewFromArray(vertices)
	setReflex(c)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		detectEars(c)
	}
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

	eliminated, err := eliminateHoles(polygon, holes)

	if err != nil {
		t.Errorf("EliminateHoles: %s", err)
	}

	checkPointArray(t, eliminated, polygonWithEliminatedHoles)
}

func TestEliminateHolesWithNoDirectlyVisible(t *testing.T) {
	// TODO
}

func TestEliminateHolesWithNoPossibleVisibleVertex(t *testing.T) {
	// TODO
}

func TestEarCut(t *testing.T) {
	result, err := EarCut(vertices, [][]Point{})
	expected := []float64{240, 115, 320, 65, 395, 170, 240, 115, 395, 170, 305, 160, 305, 160, 265, 240, 190, 100, 95, 125, 100, 215, 50, 110, 240, 115, 305, 160, 190, 100, 190, 100, 95, 125, 50, 110, 150, 30, 240, 115, 190, 100, 50, 110, 150, 30, 190, 100}

	if err != nil {
		t.Errorf("EarCut: %s", err)
	}

	t.Log(deviation(vertices, expected))
	checkFloat64Array(t, result, expected)
}

func TestEarCutSimpleShapes(t *testing.T) {
	shapes := [][]Point{
		// #0: 4 points, no reflex, results in a triangle fan
		{{0, 4}, {3, 1}, {8, 2}, {9, 5}, {4, 6}},
		// #1: diamond
		{{0, 3}, {1, 0}, {4, 1}, {3, 4}},
		// #2: square
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
		// #3: one reflex
		{{0, 6}, {0, 1}, {2, 2}, {3, 2}},
		// #4: shuriken
		{{0, 4}, {2, 2}, {2, 0}, {4, 2}, {6, 2}, {4, 4}, {4, 6}, {2, 4}},
		// #5: c letter – two reflex one after another
		{{0, 0}, {4, 0}, {4, 2}, {2, 2}, {2, 4}, {4, 4}, {4, 6}, {0, 6}},
		// #6: building
		{{1, 0}, {7, 0}, {7, 1}, {6, 1}, {6, 10}, {7, 10}, {7, 11}, {1, 11}, {1, 10}, {2, 10}, {2, 7}, {0, 7}, {0, 4}, {2, 4}, {2, 1}, {1, 1}},
	}

	for i, s := range shapes {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			res, err := EarCut(s, [][]Point{})
			if err != nil {
				t.Error(err)
			}

			real, actual, dif := deviation(s, res)

			if dif != 0 {
				t.Errorf("#%d: real area: %f; result: %f", i, real, actual)
			}
			t.Logf("#%d: %v", i, res)
		})
	}
}

func BenchmarkEarCut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EarCut(vertices, [][]Point{})
	}
}

func TestIncorrectEarCut(t *testing.T) {
	var err error
	_, err = EarCut([]Point{{0, 0}}, [][]Point{})
	if err == nil {
		t.Errorf("The code did not return error on incorrect input")
	}
}

func TestSortingByXMax(t *testing.T) {
	inners := [][]Point{{{1, 2}}, {{0, 0}}}
	sort.Sort(byMaxX(inners))
}

func TestSingleTriangleTriangulation(t *testing.T) {
	result, _ := EarCut([]Point{{0, 0}, {0, 1}, {1, 1}}, [][]Point{})
	expected := []float64{0, 0, 0, 1, 1, 1}

	checkFloat64Array(t, result, expected)
}

func TestAghA0(t *testing.T) {
	agh, _ := loadPointsFromFile("assets/agh_a0")
	for i := range agh {
		for j := range agh[i] {
			p := degreesToMeters(agh[i][j])
			agh[i][j] = Point{3 * (p.X - 2217750), 3 * (p.Y - 6457350)}
		}
	}

	result, err := EarCut(agh[0], [][]Point{}) // agh[1:]

	if err != nil {
		t.Errorf("AghA0: %s", err)
	}

	real, actual, deviation := deviation(agh[0], result)
	if deviation != 0 {
		t.Errorf("real area: %f; result: %f", real, actual)
	}
}

// **WARNING**
// Runs much longer than others (several orders of magnitude longer, can last minutes)
func TestLakeSuperior(t *testing.T) {
	t.Log("Skipping long test")
	return

	// lakeSuperior, _ := loadPointsFromFile("../../assets/lake_superior")
	// result, _ := EarCut(lakeSuperior[0], [][]Point{}) // lakeSuperior[1:]
	//
	// t.Log(result)
}

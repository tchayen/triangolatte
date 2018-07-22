package triangolatte

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"sort"
	"testing"
)

var vertices = []Point{{50, 110}, {150, 30}, {240, 115}, {320, 65}, {395, 170}, {305, 160}, {265, 240}, {190, 100}, {95, 125}, {100, 215}}

func polygonArea(data []Point) float64 {
	area := 0.0
	for i, j := 0, len(data)-1; i < len(data); i++ {
		area += data[i].X*data[j].Y - data[i].Y*data[j].X
		j = i
	}
	return math.Abs(area / 2)
}

func TestPolygonArea(t *testing.T) {
	points := []Point{{2, 2}, {11, 2}, {9, 7}, {4, 10}}
	area := polygonArea(points)

	if area != 45.5 {
		t.Error("polygonArea implementation is wrong")
	}
}

func deviation(data []Point, t []float64) (real, triangles, deviation float64) {
	trianglesArea := 0.0
	for i := 0; i < len(t); i += 6 {
		trianglesArea += math.Abs((t[0]*(t[3]-t[5]) + t[2]*(t[4]-t[1]) + t[4]*(t[1]-t[3])) / 2)
	}
	area := polygonArea(data)
	return area, trianglesArea, math.Abs(trianglesArea - area)
}

func TestDeviation(t *testing.T) {
	data := []Point{{0, 0}, {2, 0}, {2, 2}, {0, 2}}
	triangles := []float64{0, 0, 2, 0, 2, 2, 0, 2}

	t.Log(deviation(data, triangles))
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
	if cyclic(1, 5) != 1 || cyclic(4, 5) != 4 || cyclic(6, 5) != 1 || cyclic(-1, 5) != 4 || cyclic(-5, 5) != 0 || cyclic(-6, 5) != 4 {
		t.Error("cyclicList is broken")
	}
}

func TestIsReflex(t *testing.T) {
	if IsReflex(Point{0, 0}, Point{1, 1}, Point{2, 0}) != true {
		t.Error("IsReflex: false negative")
	}
	if IsReflex(Point{1, 2}, Point{0, 0}, Point{1, 0}) != false {
		t.Error("IsReflex: false positive")
	}

	// if IsReflex(Point{4, 1}, Point{1, 1}, Point{1, 4}) != false {
	// 	t.Error("Sth went wrong")
	// }
}

func BenchmarkIsReflex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsReflex(Point{0, 0}, Point{1, 1}, Point{2, 0})
	}
}

func TestIsInsideTriangle(t *testing.T) {
	case1 := IsInsideTriangle(vertices[0], vertices[8], vertices[9], vertices[7])
	case2 := IsInsideTriangle(vertices[0], vertices[1], vertices[5], vertices[7])
	if case1 == true || case2 == false {
		t.Error("IsInsideTriangle is broken")
	}
}

func BenchmarkIsInsideTriangle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsInsideTriangle(vertices[0], vertices[1], vertices[2], vertices[3])
	}
}

func TestSetReflex(t *testing.T) {
	c := NewFromArray([]Point{{0, 0}, {2, 3}, {4, 2}, {0, 7}})
	setReflex(c)

	if !(c.Front().Next().Reflex && !c.Front().Next().Next().Reflex) {
		t.Error("splitConvexAndReflex is broken")
	}
}

func TestDetectEars(t *testing.T) {
	c := NewFromArray(vertices)
	setReflex(c)
	earList := detectEars(c)

	ears := make([]Point, earList.Len())
	i, e := 0, earList.Front()

	for e != nil {
		ears[i] = e.Value.(*Element).Point
		i, e = i+1, e.Next()
	}

	third := c.Front().Next().Next().Next()
	sixth := third.Next().Next().Next()
	ninth := sixth.Next().Next().Next()
	expectedEars := []Point{third.Point, third.Next().Point, sixth.Point, ninth.Point}

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

	eliminated, _ := eliminateHoles(polygon, holes)
	checkPointArray(t, eliminated, polygonWithEliminatedHoles)
}

func TestEliminateHolesWithNoDirectlyVisible(t *testing.T) {
	// TODO
}

func TestEliminateHolesWithNoPossibleVisibleVertex(t *testing.T) {
	// TODO
}

func TestEarCut(t *testing.T) {
	result, _ := EarCut(vertices, [][]Point{})
	expected := []float64{240, 115, 320, 65, 395, 170, 240, 115, 395, 170, 305, 160, 305, 160, 265, 240, 190, 100, 95, 125, 100, 215, 50, 110, 240, 115, 305, 160, 190, 100, 190, 100, 95, 125, 50, 110, 150, 30, 240, 115, 190, 100, 50, 110, 150, 30, 190, 100}

	t.Log(deviation(vertices, expected))
	checkFloat64Array(t, result, expected)
}

func TestEarCut2(t *testing.T) {
	v := []Point{{0, 4}, {3, 1}, {8, 2}, {9, 5}, {4, 6}}
	//v := []Point{{50, 50}, {50, 200}, {200, 200}, {200, 50}}
	//v := []Point{{50, 110}, {150, 30}, {95, 125}, {100, 215}}
	//v := []Point{{10, 0}, {0, 50}, {60, 60}, {70, 10}}

	res, err := EarCut(v, [][]Point{})
	if err != nil {
		t.Error(err)
	}

	real, actual, dif := deviation(v, res)

	if dif != 0 {
		t.Errorf("real area: %f; result: %f", real, actual)
	}
	t.Log(res)
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

func loadPointsFromFile(fileName string) ([][]Point, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	polygons := make([][][]float64, 0)
	json.Unmarshal([]byte(data), &polygons)

	points := make([][]Point, len(polygons))
	for i := range polygons {
		points[i] = make([]Point, len(polygons[i]))
		for j := range polygons[i] {
			points[i][j] = Point{polygons[i][j][0], polygons[i][j][1]}
		}
	}
	return points, nil
}

func TestSingleTriangleTriangulation(t *testing.T) {
	result, _ := EarCut([]Point{{0, 0}, {0, 1}, {1, 1}}, [][]Point{})
	expected := []float64{0, 0, 0, 1, 1, 1}

	checkFloat64Array(t, result, expected)
}

func TestAghA0(t *testing.T) {
	// agh, _ := loadPointsFromFile("../../assets/agh_a0")
	// result, err := EarCut(agh[0], [][]Point{}) // agh[1:]
	//
	// t.Log(err)
	// t.Log(result)
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

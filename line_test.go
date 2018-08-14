package triangolatte

import (
	"math"
	"testing"
)

func checkArray(t *testing.T, result, expected []float64) {
	if len(result) != len(expected) {
		t.Error("Array sizes don't match")
	}

	for i := 0; i < len(result); i++ {
		if math.Abs(result[i]-expected[i]) > 0.001 {
			t.Error("Value error beyond floating point precision problem")
		}
	}
}

func TestNormal(t *testing.T) {
	line, _ := Line(Normal, []Point{{100, 300}, {100, 250}, {200, 320}, {260, 350}}, 3.0)
	expected := []float64{101.5, 250, 101.5, 300, 98.5, 300, 98.5, 300, 98.5, 250, 101.5, 250, 199.139806483455, 321.2288478807786, 99.139806483455, 251.22884788077857, 100.860193516545, 248.77115211922143, 100.860193516545, 248.77115211922143, 200.860193516545, 318.7711521192214, 199.139806483455, 321.2288478807786, 259.3291796067501, 351.34164078649985, 199.32917960675007, 321.34164078649985, 200.67082039324993, 318.65835921350015, 200.67082039324993, 318.65835921350015, 260.6708203932499, 348.65835921350015, 259.3291796067501, 351.34164078649985}

	checkArray(t, line, expected)
}

func TestMiter(t *testing.T) {
	line, _ := Line(Miter, []Point{{100, 300}, {100, 250}, {200, 320}, {260, 350}}, 3.0)
	expected := []float64{101.5, 252.88098334236005, 101.5, 300, 98.5, 300, 98.5, 300, 98.5, 247.11901665763995, 101.5, 252.88098334236005, 199.23033820382395, 321.29222008503683, 101.5, 252.88098334236005, 98.5, 247.11901665763995, 98.5, 247.11901665763995, 200.76966179617605, 318.70777991496317, 199.23033820382395, 321.29222008503683, 259.3291796067501, 351.34164078649985, 201.5, 322.8809833423601, 198.5, 317.1190166576399, 198.5, 317.1190166576399, 260.6708203932499, 348.65835921350015, 259.3291796067501, 351.34164078649985}

	checkArray(t, line, expected)
}

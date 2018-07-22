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
	line := Normal([]Point{{100, 300}, {100, 250}, {200, 320}, {260, 350}}, 3)
	expected := []float64{101, 250, 101, 300, 99, 300, 99, 300, 99, 250, 101, 250, 199.42653765563668, 320.81923192051903, 99.42653765563666, 250.81923192051903, 100.57346234436334, 249.18076807948097, 100.57346234436334, 249.18076807948097, 200.57346234436332, 319.18076807948097, 199.42653765563668, 320.81923192051903, 259.55278640450007, 350.8944271909999, 199.55278640450004, 320.8944271909999, 200.44721359549996, 319.1055728090001, 200.44721359549996, 319.1055728090001, 260.44721359549993, 349.1055728090001, 259.55278640450007, 350.8944271909999}

	checkArray(t, line, expected)
}

func TestMiter(t *testing.T) {
	line := Miter([]Point{{100, 300}, {100, 250}, {200, 320}, {260, 350}}, 3)
	expected := []float64{101, 251.92065556157337, 101, 300, 99, 300, 99, 300, 99, 248.07934443842663, 101, 251.92065556157337, 199.4868921358826, 320.8614800566912, 101, 251.92065556157337, 99, 248.07934443842663, 99, 248.07934443842663, 200.5131078641174, 319.1385199433088, 199.4868921358826, 320.8614800566912, 259.55278640450007, 350.8944271909999, 201, 321.92065556157337, 199, 318.07934443842663, 199, 318.07934443842663, 260.44721359549993, 349.1055728090001, 259.55278640450007, 350.8944271909999}

	checkArray(t, line, expected)
}

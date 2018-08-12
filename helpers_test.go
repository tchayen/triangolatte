package triangolatte

import (
	"testing"
)

func TestLoadPointsFromFile(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		_, err := LoadPointsFromFile("assets/agh_a0")

		if err != nil {
			t.Errorf("LoadPointsFromFile: %s", err)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := LoadPointsFromFile("assets/for_sure_not_there")

		if err == nil {
			t.Error("LoadPointsFromFile: returned true instead of error")
		}
	})
}

func TestDegreesToMeters(t *testing.T) {
	projected := DegreesToMeters(Point{19.57, 50.03})
	diff := projected.Sub(Point{2178522.4348243643, 6451472.9344450105})

	if diff.X > 10e-24 || diff.Y > 10.e-24 {
		t.Error("DegreesToMeters: incorrect result")
	}
}

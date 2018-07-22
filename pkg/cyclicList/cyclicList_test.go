package cyclicList

import (
	"testing"
	. "triangolatte/pkg/point"
)

func TestCyclicList_NewFromArray(t *testing.T) {
	points := []Point{{0, 1}, {1, 1}, {2, 0}}
	c := NewFromArray(points)

	for i, e := 0, c.Front(); i < c.Len(); i, e = i+1, e.Next() {
		if points[i] != e.Point {
			t.Error("Initialization from array does not work as intended")
		}
	}
}
func TestCyclicList(t *testing.T) {
	c, p1 := New(), Point{0, 1}
	p2, p3 :=Point{1, 1}, Point{2, 1}

	t.Run("push", func (t *testing.T) {
		c.Push(p1)

		if c.Len() != 1 || c.Front().Point != p1 || c.Front().Prev().Point != p1 {
			t.Error("Insertion of one point failed")
		}
	})

	t.Run("push multiple", func (t *testing.T) {
		c.Push(p2, p3)
	})

	t.Run("removal", func (t *testing.T) {
		c.Remove(c.Front().Next())

		if c.Len() != 2 || c.Front().Next().Point != p3 || c.Front().Prev().Prev().Point != p1 {
			t.Error("Removal failed")
		}
	})

	t.Run("length", func (t *testing.T) {
		if c.Len() != 2 {
			t.Error("Incorrect length")
		}
	})

	t.Run("next", func (t *testing.T) {
		if c.Front().Point != p1 {
			t.Error("Incorrect front element")
		}

		if c.Front().Prev().Point != p3 {
			t.Error("Looping does not work as intented")
		}
	})
}

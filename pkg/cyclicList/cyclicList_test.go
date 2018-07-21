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
func TestCyclicList_PushOne(t *testing.T) {
	c, p := New(), Point{1, 1}
	c.Push(p)

	if c.Len() != 1 || c.Front().Point != p || c.Front().Prev().Point != p {
		t.Error("Insertion of one point failed")
	}
}

func TestCyclicList_Push(t *testing.T) {
	c := New()
	p1, p2, p3 := Point{0, 0}, Point{1, 0}, Point{2, 0}
	c.Push(p1, p2, p3)
}

func TestCyclicList_Remove(t *testing.T) {
	c := New()
	p1, p2, p3 := Point{0, 0}, Point{1, 0}, Point{2, 0}
	c.Push(p1, p2, p3)

	c.Remove(c.Front().Next())

	if c.Len() != 2 || c.Front().Next().Point != p3 || c.Front().Prev().Prev().Point != p1 {
		t.Error("Removal failed")
	}
}

func TestCyclicList_Len(t *testing.T) {
	c := New()
	p1, p2, p3 := Point{0, 0}, Point{1, 0}, Point{2, 0}
	c.Push(p1, p2, p3)

	if c.Len() != 3 {
		t.Error("Incorrect length")
	}
}

func TestElement_Next(t *testing.T) {
	c := New()
	p1, p2, p3 := Point{0, 1}, Point{1, 1}, Point{2, 1}
	c.Push(p1, p2, p3)

	if c.Front().Point != p1 {
		t.Error("Incorrect front element")
	}

	if c.Front().Prev().Point != p3 {
		t.Error("Looping does not work as intented")
	}
}

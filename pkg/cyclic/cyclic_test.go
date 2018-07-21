package cyclic

import (
    "testing"
    . "triangolatte/pkg/point"
)

func TestElement_Next(t *testing.T) {
    c := New()
    p1, p2, p3 := Point{0, 0}, Point{1, 0}, Point{2, 0}
    c.Push(p1, p2, p3)

    if c.Len() != 3 {
        t.Error("Incorrect length")
    }

    if c.Front().Point != p1 {
        t.Error("Incorrect front element")
    }

    if c.Front().Prev().Point != p3 {
        t.Error("Looping does not work as intented")
    }
}

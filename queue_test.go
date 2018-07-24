package triangolatte

import "testing"

func TestQueue(t *testing.T) {
	var q *Queue

	t.Run("NewQueue", func(t *testing.T) {
		q = NewQueue(3)
	})

	e1, e2 := Element{Point: Point{1, 2}}, Element{Point: Point{3, 4}}

	t.Run("Push", func(t *testing.T) {
		q.Push(&e1)
		q.Push(&e2)
	})

	t.Run("Pop", func(t *testing.T) {
		popped := q.Pop()

		if popped != &e1 {
			t.Error("Pop returned wrong element")
		}
	})

	e3 := Element{Point: Point{5, 6}}

	t.Run("Push again", func(t *testing.T) {
		q.Push(&e3)
	})

	e4 := Element{Point: Point{7, 8}}

	t.Run("Push to the capacity with index wrap", func(t *testing.T) {
		q.Push(&e4)
	})
}

package triangolatte

import "testing"

func checkQueue(t *testing.T, q *Queue, expected []*Element) {
	for i := 0; i < q.capacity; i++ {
		if q.elements[i] != expected[i] {
			t.Error("Error")
		}
	}
}

func TestQueue(t *testing.T) {
	var q *Queue

	t.Run("new queue", func(t *testing.T) {
		q = NewQueue(3)
	})

	e1, e2 := Element{Point: Point{1, 2}}, Element{Point: Point{3, 4}}

	t.Run("push", func(t *testing.T) {
		q.Push(&e1)
		q.Push(&e2)
		expected := []*Element{&e1, &e2, nil}
		checkQueue(t, q, expected)
	})

	t.Run("pop", func(t *testing.T) {
		q.Pop()
		expected := []*Element{nil, &e2, nil}
		checkQueue(t, q, expected)
	})

	e3 := Element{Point: Point{5, 6}}

	t.Run("push again", func(t *testing.T) {
		q.Push(&e3)
		expected := []*Element{nil, &e2, &e3}
		checkQueue(t, q, expected)
	})

	e4 := Element{Point: Point{7, 8}}

	t.Run("push to the capacity with index wrap", func(t *testing.T) {
		q.Push(&e4)
		expected := []*Element{&e4, &e2, &e3}
		checkQueue(t, q, expected)
	})

	t.Run("remove", func(t *testing.T) {
		q.Remove(1)
		expected := []*Element{nil, &e4, &e3}
		checkQueue(t, q, expected)
	})
}

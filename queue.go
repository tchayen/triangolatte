package triangolatte

import "errors"

type Queue struct {
	points []*Element
	start, end, size, capacity int
}

func (q *Queue) Init(n int) *Queue {
	q.points = make([]*Element, n)
	q.start = 0
	q.end = n - 1
	q.size = 0
	q.capacity = n
	return q
}

func NewQueue(n int) *Queue {
	return new(Queue).Init(n)
}

func (q *Queue) Pop() *Element {
	e := q.points[q.start]
	q.points[q.start] = nil
	q.start = (q.start + 1) % q.capacity
	q.size--
	return e
}

func (q *Queue) Push(e *Element) error {
	if q.size == q.capacity {
		return errors.New("cannot push element to queue – maximal capacity reached")
	}

	q.end = (q.end + 1) % q.capacity
	q.points[q.end] = e
	q.size++
	return nil
}

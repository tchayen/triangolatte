package triangolatte

type Queue struct {
	elements   []*Element
	start, end int
	len        int
	capacity   int
}

func (q *Queue) Init(n int) *Queue {
	q.elements = make([]*Element, n)
	q.start = 0
	q.end = 0
	q.len = 0
	q.capacity = n
	return q
}

func NewQueue(n int) *Queue {
	return new(Queue).Init(n)
}

func (q *Queue) Pop() *Element {
	if q.len == 0 {
		panic("Cannot pop from empty queue")
	}

	e := q.elements[q.start]
	q.elements[q.start] = nil
	q.len--

	if q.len > 0 {
		q.start = (q.start + 1) % q.capacity
	} else {
		q.start = 0
		q.end = 0
	}

	return e
}

func (q *Queue) Push(e *Element) int {
	if q.len == q.capacity {
		panic("Cannot push to full queue")
	}

	if q.len > 0 {
		q.end = (q.end + 1) % q.capacity
	}

	q.elements[q.end] = e
	q.len++
	return q.end
}

func (q *Queue) Remove(i int) {
	if q.len == 0 {
		panic("Cannot remove from empty queue")
	}

	q.elements[i] = q.elements[q.end]
	q.elements[q.end] = nil
	q.len--

	if q.len > 0 {
		q.end = ((q.end-1)%q.capacity + q.capacity) % q.capacity
	} else {
		q.start = 0
		q.end = 0
	}
}

func (q *Queue) Len() int {
	return q.len
}

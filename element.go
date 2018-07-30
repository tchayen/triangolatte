package triangolatte

type Element struct {
	Prev, Next *Element
	Point      Point
}

func Insert(p Point, e *Element) *Element {
	new := Element{Point: p}

	if e != nil {
		new.Next = e.Next
		new.Prev = e
		e.Next.Prev = &new
		e.Next = &new
	} else {
		new.Prev = &new
		new.Next = &new
	}
	return &new
}

func (e *Element) Remove() {
	e.Next.Prev = e.Prev
	e.Prev.Next = e.Next
}

package triangolatte

import (
	"container/list"
)

// Implements double linked cyclic lists.
//
// NOTE: in one element list e == e.next
//
// Parameters
//
// - removal O(1)
//
// - addition O(1)
//
// - length O(1)
//
// Initialize using
//      c := cyclic.New()
//
// Iterate over the list with
//      for i, e := 0, list.Front(); i < list.Len(); i, e = i + 1, e.Next() {
//          // do something with e.Point, e.Reflex, e.Ear...
//      }
//
type CyclicList struct {
	root Element
	len  int
}

type Element struct {
	// The list to which this element belongs.
	List *CyclicList

	// Next and previous elements.
	prev, next *Element

	// Value of the element.
	Point Point

	// Its properties.
	Reflex bool
	Ear    *list.Element
}

// Init sets up fresh CyclicList.
func (c *CyclicList) Init() *CyclicList {
	c.root.next = &c.root
	c.root.prev = &c.root
	c.root.List = c
	c.len = 0
	return c
}

// New call new(CyclicList) and initializes it.
func New() *CyclicList {
	return new(CyclicList).Init()
}

// NewFromArray takes []Point and puts it inside CyclicList.
func NewFromArray(points []Point) *CyclicList {
	c := New()
	after := &c.root
	for _, p := range points {
		after = c.InsertAfter(p, after)
	}
	return c
}

// Len returns length of the CyclicList.
func (c *CyclicList) Len() int {
	return c.len
}

// Front returns first element (head) of the CyclicList.
func (c *CyclicList) Front() *Element {
	if c.len == 0 {
		return nil
	}
	return c.root.next
}

// InsertAfter puts given Point after given *Element (as its next element and
// previous to the given element's successor)
func (c *CyclicList) InsertAfter(p Point, e *Element) *Element {
	new := Element{Point: p, prev: e, next: e.next, List: e.List}
	e.next.prev = &new
	e.next = &new
	c.len++
	return &new
}

// Push inserts variable number of points to the CyclicList.
func (c *CyclicList) Push(points ...Point) {
	after := c.root.prev
	for _, p := range points {
		after = c.InsertAfter(p, after)
	}
}

// Remove detaches given *Element from the CyclicList and resets its pointers.
func (c *CyclicList) Remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev

	// Avoid memory leaks.
	e.next = nil
	e.prev = nil
	e.List = nil

	c.len--

	return e
}

// Next returns next element on the list (or first if given element was
// the last one).
func (e *Element) Next() *Element {
	if e.List == nil {
		return nil
	}

	if e.next == &e.List.root {
		return e.List.root.next
	}

	return e.next
}

// Prev returns previous element on the list (or last one if given element
// was the first one).
func (e *Element) Prev() *Element {
	if e.List == nil {
		return nil
	}

	if e.prev == &e.List.root {
		return e.List.root.prev
	}

	return e.prev
}

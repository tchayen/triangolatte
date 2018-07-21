// Package implements double linked cyclic lists.
//
// **NOTE:** _in one element list `e == e.next`_
//
// Initialize using
//      c := cyclic.New()
//
// Iterate over the list with
//      for i, e := 0, list.Front(); i < list.Len(); i, e = i + 1, e.Next() {
//          // do something with e.Point, e.Reflex, e.Ear...
//      }
//

package cyclic

import . "triangolatte/pkg/point"

type Cyclic struct {
    root Element
    len int
}

type Element struct {
    // The list to which this element belongs.
    list *Cyclic

    // Next and previous elements.
    prev, next *Element

    // Value of the element.
    Point Point

    // Its properties.
    Reflex bool
    Ear bool
}

func (c *Cyclic) Init() *Cyclic {
    c.root.next = &c.root
    c.root.prev = &c.root
    c.len = 0
    return c
}

func New() *Cyclic {
    return new(Cyclic).Init()
}

func (c *Cyclic) First() *Element {
    return &c.root
}

func (c *Cyclic) Len() int {
    return c.len
}

func (c *Cyclic) Front() *Element {
    if c.len == 0 {
        return nil
    }
    return c.root.next
}

func (c *Cyclic) InsertAfter(p Point, e *Element) *Element {
    new := Element{Point: p, prev: e, next: e.next, list: e.list}
    e.next.prev = &new
    e.next = &new
    c.len++
    return &new
}

func (c *Cyclic) Push(points ...Point) {
    for _, p := range points {
        c.InsertAfter(p, c.root.prev)
    }
}

func (e *Element) Next() *Element {
    if e.list == nil {
        return nil
    }

    if e.next == &e.list.root {
        return e.list.root.next
    }

    return e.next
}

func (e *Element) Prev() *Element {
    if e.list == nil {
        return nil
    }

    if e.prev == &e.list.root {
        return e.list.root.prev
    }

    return e.prev
}

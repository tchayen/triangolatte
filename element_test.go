package triangolatte

import "testing"

func TestInsert(t *testing.T) {
	t.Run("from scratch", func(t *testing.T) {
		e := Insert(Point{1, 1}, nil)
		if e.Next != e || e.Prev != e {
			t.Error("Wrong creation of single element list")
		}
	})

	t.Run("two element", func(t *testing.T) {
		e1 := Insert(Point{1, 1}, nil)
		e2 := Insert(Point{2, 2}, e1)

		if e1.Next != e2 || e2.Next != e1 || e2.Prev != e1 {
			t.Error("Wrong two element list")
		}
	})
}

func TestElement_Remove(t *testing.T) {
	t.Run("removal", func(t *testing.T) {
		e1 := Insert(Point{1, 1}, nil)
		e2 := Insert(Point{2, 2}, e1)
		e3 := Insert(Point{3, 3}, e2)

		if e1.Next != e2 || e2.Next != e3 {
			t.Error("Wrong insert")
		}

		e2.Remove()

		if e1.Next != e3 || e3.Prev != e1 {
			t.Error("Removal did not connect outer nodes")
		}

		if e2.Prev != e1 || e2.Next != e3 {
			t.Error("Removal did not preserve connections")
		}
	})

	t.Run("remove edge", func(t *testing.T) {
		e1 := Insert(Point{1, 1}, nil)
		e2 := Insert(Point{2, 2}, e1)
		e3 := Insert(Point{3, 3}, e2)

		e3.Remove()

		if e2.Next != e1 || e1.Prev != e2 {
			t.Error("Edge removal makes incorrect connection")
		}
	})
}

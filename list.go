package main

type TDequeue struct {
	head *T
	tail *T
}

func (td *TDequeue) push(t *T) {
	if td.head != nil {
		td.head.prev = t
	} else {
		td.tail = t
	}
	td.head, t.next = t, td.head
}

func (td *TDequeue) pop() *T {
	if td.head != nil {
		if td.head == td.tail {
			o := td.head
			td.head, td.tail = nil, nil
			return o
		}
		o := td.head
		td.head = td.head.next
		td.head.prev = nil
		return o
	}
	return nil
}

func (td *TDequeue) pushBack(t *T) {
	if td.tail != nil {
		td.tail.next = t
		t.prev = td.tail
		td.tail = t
		return
	}
	td.head = t
	td.tail = t
}

func (td *TDequeue) popTail() *T {
	if td.tail != nil {
		o := td.tail
		td.tail = o.prev
		if o.prev != nil {
			td.tail.next = nil
		}
		return o
	}
	return nil
}

// append appends a list (in order) to the
// end of the current list.
func (td *TDequeue) append(d TDequeue) {
	td.tail.next = d.head
	d.head.prev = td.tail
	if d.tail != nil {
		td.tail = d.tail
	}
}

// prepend puts a list in front of the current list.
func (td *TDequeue) prepend(d TDequeue) {
	d.tail.next = td.head
	td.head.prev = d.tail
	if d.head != nil {
		td.head = d.head
	}
}

// filter removes elements from the current list and puts
// them into the returned list. the ordering of elements is
// preserved.
func (td *TDequeue) filter(match func(*T) bool) TDequeue {
	out := TDequeue{}
	hd := td.head
	var cur *T
	for hd != nil {
		cur = hd.next
		if match(hd) {
			// unlink
			if hd == td.head {
				td.pop()
			} else if hd == td.tail {
				td.popTail()
			} else {
				hd.prev.next, hd.next.prev = hd.next, hd.prev
				hd.next = nil
				hd.prev = nil
			}
			out.pushBack(hd)
		}
		hd = cur
	}
	return out
}

//
func (td *TDequeue) count() int {
	i := 0
	for t := td.head; t != nil; t = t.next {
		i++
	}
	return i
}

// 'at' returns the 'n'th element in the list,
// or 'nil' if the index is out of range.
func (td *TDequeue) at(n int) *T {
	var o *T
	for o = td.head; o != nil && n > 0; n-- {
		o = o.next
	}
	return o
}

// insertAt sets 't' to be the 'n'th element.
// if 'n' is out of bounds, 't' is appended to
// the tail of the queue.
func (td *TDequeue) insertAt(n int, t *T) {
	if n == 0 {
		td.push(t)
		return
	}
	el := td.at(n - 1)
	if el == nil || el == td.tail {
		td.pushBack(t)
		return
	}
	t.prev, t.next, el.next = el, el.next, t
}

// removeAt removes the element at index 'n'.
// if 'n' is out of bounds, the return value is 'nil'.
func (td *TDequeue) removeAt(n int) *T {
	if n == 0 {
		return td.pop()
	}
	el := td.at(n)
	if el == nil {
		return nil
	} else if el == td.tail {
		return td.popTail()
	}
	el.prev.next, el.next.prev = el.next, el.prev
	return el
}

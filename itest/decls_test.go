package itest

import (
	"testing"
)

// returns the index at which the nodes are not equal,
// or (na.count(), true)
func assertEqual(na *nodeList, nb *nodeList) (int, bool) {
	elA := na.head
	elB := nb.head
	i := 0
	for elA != nil && elB != nil {
		if elA.value != elB.value {
			return i, false
		}
		elA = elA.next
		elB = elB.next
		i++
	}
	// both els should be nil
	return i, (elA == elB)
}

func assertCount(nl *nodeList, n int, t *testing.T) {
	if nl.count() != n {
		t.Errorf("expected list.count() to be %d; got %d", n, nl.count())
	}
}

func TestBasic(t *testing.T) {
	list := nodeList{}

	first := &node{value: 1}

	list.pushBack(first)
	assertCount(&list, 1, t)

	list.pushBack(&node{value: 2})
	assertCount(&list, 2, t)

	top := list.pop()
	if top.value != 1 {
		t.Errorf("expected top node to be 1; got %d", top.value)
	}
	assertCount(&list, 1, t)

	top = list.pop()
	if top.value != 2 {
		t.Errorf("expected top node to be 2; got %d", top.value)
	}
	assertCount(&list, 0, t)

	// invert the list
	list.push(first)
	list.push(top)

	assertCount(&list, 2, t)

	if list.at(0).value != 2 {
		t.Errorf("expected list.at(0).value to be 2; got %d", list.at(0).value)
	}
	if list.at(1).value != 1 {
		t.Errorf("expected list.at(1).value to be 1; got %d", list.at(1).value)
	}
	if list.at(2) != nil {
		t.Errorf("expected list.at(2) to be nil; got %+#v", list.at(2))
	}

	// make list {2, 100, 1}
	list.insertAt(1, &node{value: 100})
	assertCount(&list, 3, t)
	if list.at(1).value != 100 {
		t.Errorf("expected list.at(1).value to be 100; got %d", list.at(1).value)
	}
	if list.at(-1) != list.tail {
		t.Errorf("expected list.at(-1) (%+#v) to be list.tail (%+#v)", list.at(-1), list.tail)
	}
	if list.at(-2) != list.at(1) {
		t.Errorf("expected list.at(-2) (%+#v) to be list.at(-1) (%+#v)", list.at(-2), list.at(1))
	}

	// make list {2, 1}
	hundred := list.removeAt(1)
	assertCount(&list, 2, t)
	if hundred.value != 100 {
		t.Errorf("expected list.removeAt(1) to return 100; got %d", hundred.value)
	}

	list.popTail()
	list.popTail()
	assertCount(&list, 0, t)
}

func TestFilter(t *testing.T) {
	list := nodeList{}

	// ordered list from 0 to 100
	for i := 0; i < 100; i++ {
		list.pushBack(&node{value: i})
	}
	assertCount(&list, 100, t)

	evens := list.filter(func(n *node) bool { return n.value&1 == 0 })

	list.walk(func(n *node, idx int) {
		if n.value&1 == 0 {
			t.Errorf("element %d of original list (%d) is even", idx, n.value)
		}
	})

	evens.walk(func(n *node, idx int) {
		if n.value&1 != 0 {
			t.Errorf("element %d of evens (%d) is odd", idx, n.value)
		}
	})

	if list.count() != evens.count() {
		t.Errorf("expected evens == odds; got %d evens and %d odds", evens.count(), list.count())
	}

	list.prepend(evens)
	assertCount(&list, 100, t)

	// flip
	list.append(list.filter(func(n *node) bool { return n.value%3 == 0 }))
	assertCount(&list, 100, t)
}

func BenchmarkPush1000(b *testing.B) {
	nodes := make([]node, 1000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := nodeList{}
		for n := 0; n < 1000; n++ {
			list.push(&nodes[n])
		}
	}
}

func BenchmarkPushPop1000(b *testing.B) {
	nodes := make([]node, 1000)
	b.ReportAllocs()
	b.ResetTimer()
	list := nodeList{}
	for i := 0; i < b.N; i++ {
		for n := 0; n < 1000; n++ {
			list.push(&nodes[n])
		}
		for n := 0; n < 1000; n++ {
			list.pop()
		}
	}
}

func BenchmarkFilter1000(b *testing.B) {
	list := nodeList{}
	nodes := make([]node, 1000)
	for i := 0; i < 1000; i++ {
		nodes[i].value = i
		list.pushBack(&nodes[i])
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.append(list.filter(func(n *node) bool { return n.value&1 == 0 }))
	}
}

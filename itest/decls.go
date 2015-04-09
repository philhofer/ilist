package itest

//go:generate ilist -type=node

type node struct {
	value	int
	next	*node
	prev	*node
}

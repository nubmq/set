package set

import "fmt"

// Previous Node and Set implementations remain the same...
// [Previous code from Color definition through deleteFixup remains unchanged]

// Iterator represents a bidirectional iterator for the set
type Iterator struct {
	node    *Node
	set     *Set
	reverse bool
}

// Begin returns an iterator to the smallest element
func (s *Set) Begin() *Iterator {
	if s.root == nil {
		return &Iterator{set: s}
	}
	return &Iterator{
		node: s.minimum(s.root),
		set:  s,
	}
}

// End returns an iterator past the largest element
func (s *Set) End() *Iterator {
	return &Iterator{set: s}
}

// RBegin returns a reverse iterator to the largest element
func (s *Set) RBegin() *Iterator {
	if s.root == nil {
		return &Iterator{set: s, reverse: true}
	}
	return &Iterator{
		node:    s.maximum(s.root),
		set:     s,
		reverse: true,
	}
}

// REnd returns a reverse iterator before the smallest element
func (s *Set) REnd() *Iterator {
	return &Iterator{set: s, reverse: true}
}

// Value returns the current element
func (it *Iterator) Value() interface{} {
	if it.node == nil {
		return nil
	}
	return it.node.key
}

// Valid returns true if the iterator points to a valid element
func (it *Iterator) Valid() bool {
	return it.node != nil
}

// Next moves to the next element
// Returns false if we've reached the end
func (it *Iterator) Next() bool {
	if it.node == nil {
		return false
	}
	
	if it.reverse {
		it.node = it.set.predecessor(it.node)
	} else {
		it.node = it.set.successor(it.node)
	}
	
	return it.node != nil
}

// Prev moves to the previous element
// Returns false if we've reached the beginning
func (it *Iterator) Prev() bool {
	if it.node == nil {
		if it.reverse {
			it.node = it.set.minimum(it.set.root)
		} else {
			it.node = it.set.maximum(it.set.root)
		}
		return it.node != nil
	}
	
	if it.reverse {
		it.node = it.set.successor(it.node)
	} else {
		it.node = it.set.predecessor(it.node)
	}
	
	return it.node != nil
}

// Helper function to find maximum element
func (s *Set) maximum(x *Node) *Node {
	for x.right != nil {
		x = x.right
	}
	return x
}

// Helper function to find predecessor
func (s *Set) predecessor(x *Node) *Node {
	if x.left != nil {
		return s.maximum(x.left)
	}
	y := x.parent
	for y != nil && x == y.left {
		x = y
		y = y.parent
	}
	return y
}

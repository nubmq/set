// Package set provides a Red-Black Tree based Set implementation
package set

// Color represents the color of a node in the Red-Black tree
type Color bool

const (
	Black Color = true
	Red   Color = false
)

// Node represents a node in the Red-Black tree
type Node struct {
	key                     interface{}
	color                   Color
	left, right, parent     *Node
}

// Set represents the Red-Black tree based set
type Set struct {
	root    *Node
	size    int
	compare func(interface{}, interface{}) int
}

// Iterator represents a bidirectional iterator for the set
type Iterator struct {
	node    *Node
	set     *Set
	reverse bool
}

// NewSet creates a new set with a custom comparator
func NewSet(compare func(interface{}, interface{}) int) *Set {
	return &Set{
		compare: compare,
	}
}

// Size returns the number of elements in the set
func (s *Set) Size() int {
	return s.size
}

// Clear removes all elements from the set
func (s *Set) Clear() {
	s.root = nil
	s.size = 0
}

// IsEmpty returns true if the set has no elements
func (s *Set) IsEmpty() bool {
	return s.size == 0
}

// Insert adds a new element to the set
func (s *Set) Insert(key interface{}) bool {
	if s.root == nil {
		s.root = &Node{
			key:     key,
			color:   Black,
		}
		s.size++
		return true
	}

	node := s.root
	var parent *Node

	for node != nil {
		parent = node
		cmp := s.compare(key, node.key)
		if cmp == 0 {
			return false // Key already exists
		} else if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}

	newNode := &Node{
		key:     key,
		color:   Red,
		parent:  parent,
	}

	if s.compare(key, parent.key) < 0 {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	s.size++
	s.insertFixup(newNode)
	return true
}

// Contains checks if an element exists in the set
func (s *Set) Contains(key interface{}) bool {
	node := s.root
	for node != nil {
		cmp := s.compare(key, node.key)
		if cmp == 0 {
			return true
		} else if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}
	return false
}

// Remove removes an element from the set
func (s *Set) Remove(key interface{}) bool {
	node := s.root
	for node != nil {
		cmp := s.compare(key, node.key)
		if cmp == 0 {
			s.delete(node)
			s.size--
			return true
		} else if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}
	return false
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

// Iterator methods
func (it *Iterator) Value() interface{} {
	if it.node == nil {
		return nil
	}
	return it.node.key
}

func (it *Iterator) Valid() bool {
	return it.node != nil
}

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

// Internal helper functions
func (s *Set) leftRotate(x *Node) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		s.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (s *Set) rightRotate(x *Node) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		s.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

func (s *Set) insertFixup(z *Node) {
	for z.parent != nil && z.parent.color == Red {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y != nil && y.color == Red {
				z.parent.color = Black
				y.color = Black
				z.parent.parent.color = Red
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					s.leftRotate(z)
				}
				z.parent.color = Black
				z.parent.parent.color = Red
				s.rightRotate(z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y != nil && y.color == Red {
				z.parent.color = Black
				y.color = Black
				z.parent.parent.color = Red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					s.rightRotate(z)
				}
				z.parent.color = Black
				z.parent.parent.color = Red
				s.leftRotate(z.parent.parent)
			}
		}
	}
	s.root.color = Black
}

func (s *Set) minimum(x *Node) *Node {
	for x.left != nil {
		x = x.left
	}
	return x
}

func (s *Set) maximum(x *Node) *Node {
	for x.right != nil {
		x = x.right
	}
	return x
}

func (s *Set) successor(x *Node) *Node {
	if x.right != nil {
		return s.minimum(x.right)
	}
	y := x.parent
	for y != nil && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

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

func (s *Set) delete(z *Node) {
	var x, y *Node
	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = s.successor(z)
	}

	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	if x != nil {
		x.parent = y.parent
	}

	if y.parent == nil {
		s.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != z {
		z.key = y.key
	}

	if y.color == Black {
		s.deleteFixup(x, y.parent)
	}
}

func (s *Set) deleteFixup(x *Node, parent *Node) {
	for x != s.root && (x == nil || x.color == Black) {
		if x == parent.left {
			w := parent.right
			if w.color == Red {
				w.color = Black
				parent.color = Red
				s.leftRotate(parent)
				w = parent.right
			}
			if (w.left == nil || w.left.color == Black) &&
				(w.right == nil || w.right.color == Black) {
				w.color = Red
				x = parent
				parent = x.parent
			} else {
				if w.right == nil || w.right.color == Black {
					if w.left != nil {
						w.left.color = Black
					}
					w.color = Red
					s.rightRotate(w)
					w = parent.right
				}
				w.color = parent.color
				parent.color = Black
				if w.right != nil {
					w.right.color = Black
				}
				s.leftRotate(parent)
				x = s.root
			}
		} else {
			w := parent.left
			if w.color == Red {
				w.color = Black
				parent.color = Red
				s.rightRotate(parent)
				w = parent.left
			}
			if (w.right == nil || w.right.color == Black) &&
				(w.left == nil || w.left.color == Black) {
				w.color = Red
				x = parent
				parent = x.parent
			} else {
				if w.left == nil || w.left.color == Black {
					if w.right != nil {
						w.right.color = Black
					}
					w.color = Red
					s.leftRotate(w)
					w = parent.left
				}
				w.color = parent.color
				parent.color = Black
				if w.left != nil {
					w.left.color = Black
				}
				s.rightRotate(parent)
				x = s.root
			}
		}
	}
	if x != nil {
		x.color = Black
	}
}

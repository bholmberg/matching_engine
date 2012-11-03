package trade

import ()

type Tree struct {
	size int
	root *Node
}

func NewTree() *Tree {
	return &Tree{}
}

func (b *Tree) Push(n *Node) {
	if b.root != nil {
		b.root.push(n)
	} else {
		b.root = n
		n.pp = &b.root
	}
}

func (b *Tree) PeekMin() *Node {
	return b.root.peekMin()
}

func (b *Tree) PopMin() *Node {
	if b.root != nil {
		b.size--
		return b.root.popMin()
	}
	return nil
}

func (b *Tree) PeekMax() *Node {
	return b.root.peekMax()
}

func (b *Tree) PopMax() *Node {
	if b.root != nil {
		b.size--
		return b.root.popMax()
	}
	return nil
}

func (b *Tree) Pop(val int64) *Node {
	n := b.root.popVal(val)
	if n != nil {
		b.size--
	}
	return n
}

func (b *Tree) Size() int {
	return b.size
}

type Node struct {
	pp    **Node
	left  *Node
	right *Node
	next  *Node
	prev  *Node
	size  int
	val   int64
	O     *Order
}

func initNode(val int64, o *Order, n *Node) {
	*n = Node{val: val, O: o}
	n.next = n
	n.prev = n
}

func (n *Node) peekMax() *Node {
	if n == nil {
		return nil
	}
	if n.right == nil {
		return n
	}
	return n.right.peekMax()
}

func (n *Node) popMax() *Node {
	if n == nil {
		return nil
	}
	m := n.peekMax()
	return m.pop()
}

func (n *Node) peekMin() *Node {
	if n == nil {
		return nil
	}
	if n.left == nil {
		return n
	}
	return n.left.peekMin()
}

func (n *Node) popMin() *Node {
	if n == nil {
		return nil
	}
	m := n.peekMin()
	return m.pop()
}

func (n *Node) popVal(val int64) *Node {
	switch {
	case n == nil:
		return nil
	case val == n.val:
		return n.pop()
	case val < n.val:
		return n.left.popVal(val)
	case val > n.val:
		return n.right.popVal(val)
	}
	panic("unreachable")
}

func (n *Node) push(in *Node) {
	switch {
	case in.val == n.val:
		n.insert(in)
	case in.val < n.val:
		if n.left == nil {
			in.pp = &n
			n.left = in
		} else {
			n.left.push(in)
		}
	case in.val > n.val:
		if n.right == nil {
			in.pp = &n
			n.right = in
		} else {
			n.right.push(in)
		}
	}
}

func (n *Node) insert(in *Node) {
	last := n.next
	last.prev = in
	in.next = last
	in.prev = n
	n.next = n
	n.size++
}

func (n *Node) pop() *Node {
	n.prev.next = n.next
	n.next.prev = n.prev
	nn := n.next
	nn.size = n.size - 1
	n.next = nil
	n.prev = nil
	n.size = 1
	if n != nn {
		swap(n, nn)
	} else {
		detatch(n)
	}
	return n
}

func swap(n *Node, nn *Node) {
	nn.pp = n.pp
	nn.left = n.left
	nn.right = n.right
	*nn.pp = nn
	nn.left.pp = &nn
	nn.right.pp = &nn
}

func detatch(n *Node) {
	switch {
	case n.right == nil && n.left == nil:
		*n.pp = nil
	case n.right == nil:
		*n.pp = n.left
		n.left.pp = n.pp
	case n.left == nil:
		*n.pp = n.right
		n.right.pp = n.pp
	default:
		nn := n.left.detatchMax()
		swap(n, nn)
	}
	n.pp = nil
	n.left = nil
	n.right = nil
}

func (n *Node) detatchMax() *Node {
	m := n.peekMax()
	detatch(m)
	return m
}
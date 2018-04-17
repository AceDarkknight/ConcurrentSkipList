package ConcurrentSkipList

type Node struct {
	index     uint64
	value     interface{}
	nextNodes []*Node
}

// newNode will create a node using in this package but not external package.
func newNode(index uint64, value interface{}, level int) *Node {
	return &Node{
		index:     index,
		value:     value,
		nextNodes: make([]*Node, level, level),
	}
}

// Index will return the node's index.
func (n *Node) Index() uint64 {
	return n.index
}

// Value will return the node's value.
func (n *Node) Value() interface{} {
	return n.value
}

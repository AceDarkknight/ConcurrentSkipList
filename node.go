package ConcurrentSkipList

type Node struct {
	index     uint64
	value     interface{}
	nextNodes []*Node
}

func newNode(index uint64, value interface{}, level int) *Node {
	if level <= 0 || level > MAX_LEVEL {
		level = MAX_LEVEL
	}

	return &Node{
		index:     index,
		value:     value,
		nextNodes: make([]*Node, level),
	}
}

func (n *Node) Index() uint64 {
	return n.index
}

func (n *Node) Value() interface{} {
	return n.value
}

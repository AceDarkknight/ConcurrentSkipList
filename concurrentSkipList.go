package ConcurrentSkipList

import (
	"github.com/OneOfOne/xxhash"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
)

// Comes from redis's implementation.
// Also you can see more detail in William Pugh's paper <Skip Lists: A Probabilistic Alternative to Balanced Trees>.
// The paper is in ftp://ftp.cs.umd.edu/pub/skipLists/skiplists.pdf
const (
	MAX_LEVEL   = 32
	PROBABILITY = 0.25
	SHARDS      = 32
)

var shardIndexes = make([]uint64, SHARDS)

func init() {
	var step uint64 = 1 << 59 // 2^64/32
	shardIndexes[0] = step
	for i := SHARDS - 1; i > 0; i-- {
		var t uint64 = math.MaxUint64
		for j := SHARDS - 1; j < i; i++ {
			t -= step
		}

		shardIndexes[i] = t
	}
}

type Node struct {
	index        uint64
	value        interface{}
	previousNode *Node
	nextNodes    []*Node
}

func NewNode(index uint64, value interface{}, level int) *Node {
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

func (n *Node) Next() *Node {
	return n.nextNodes[0]
}

func (n *Node) Previous() *Node {
	return n.previousNode
}

type ConcurrentSkipList struct {
	skipLists []*skipList
	length    int32
}

type skipList struct {
	level  int
	length int32
	head   *Node
	tail   *Node
	max    *Node
	mutex  sync.RWMutex
}

func NewConcurrentSkipList(level int) *ConcurrentSkipList {
	if level <= 0 || level > MAX_LEVEL {
		level = MAX_LEVEL
	}

	skipLists := make([]*skipList, SHARDS)
	for i := 0; i < SHARDS; i++ {
		head := NewNode(0, nil, level)
		tail := NewNode(math.MaxUint64, nil, level)
		for i := 0; i < len(head.nextNodes); i++ {
			head.nextNodes[i] = tail
		}

		tail.previousNode = head
		head.previousNode = nil

		skipLists[i] = &skipList{
			level:  level,
			length: 0,
			head:   head,
			tail:   tail,
		}
	}

	return &ConcurrentSkipList{
		skipLists: skipLists,
	}
}

// Level will return the level of skip list.
func (s *ConcurrentSkipList) Level() int {
	return s.skipLists[0].level
}

// Length will return the length of skip list.
func (s *ConcurrentSkipList) Length() int32 {
	var length int32 = 0
	for _, sl := range s.skipLists {
		length += atomic.LoadInt32(&sl.length)
	}

	return length
}

func (s *ConcurrentSkipList) Search(index uint64) *Node {
	sl := s.skipLists[getShardIndex(index)]
	if atomic.LoadInt32(&sl.length) == 0 {
		return nil
	}

	currentNode := sl.head
	for l := sl.level - 1; l >= 0; l-- {
		for currentNode.nextNodes[l] != sl.tail && currentNode.index < index {
			currentNode = currentNode.nextNodes[l]
		}
	}

	if currentNode == sl.tail || currentNode.index > index {
		return nil
	} else if currentNode.index == index {
		return currentNode
	} else {
		return nil
	}
}

func getShardIndex(index uint64) int {
	result := -1
	for i, t := range shardIndexes {
		if index <= t {
			result = i
			break
		}
	}

	return result
}

// randomLevel will generate and random level that level > 0 and level < skip list's level
// This comes from redis's implementation.
func (s *skipList) randomLevel() int {
	level := 1
	for rand.Float64() < PROBABILITY && level < s.level {
		level++
	}

	return level
}

// Hash will calculate the input's hash value using xxHash algorithm.
// It can be used to calculate the index of skip list.
// See more detail in https://cyan4973.github.io/xxHash/
func Hash(input []byte) uint64 {
	h := xxhash.New64()
	h.Write(input)
	return h.Sum64()
}

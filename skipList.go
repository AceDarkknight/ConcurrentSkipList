package ConcurrentSkipList

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

type skipList struct {
	level  int
	length int32
	head   *Node
	tail   *Node
	mutex  sync.RWMutex
}

func newSkipList(level int) *skipList {
	head := newNode(0, nil, level)
	var tail *Node = nil
	for i := 0; i < len(head.nextNodes); i++ {
		head.nextNodes[i] = tail
	}

	return &skipList{
		level:  level,
		length: 0,
		head:   head,
		tail:   tail,
	}
}

// searchWithPreviousNode will search given index in skip list.
// The first return value represents the previous nodes need to update when call Insert function.
// The second return value represents the value with given index or the closet value whose index is larger than given index.
func (s *skipList) searchWithPreviousNodes(index uint64) ([]*Node, *Node) {
	// Store all previous value whose index is less than index and whose next value's index is larger than index.
	previousNodes := make([]*Node, s.level)

	// fmt.Printf("start doSearch:%v\n", index)
	currentNode := s.head

	// Lock and unlock.
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Iterate from top level to bottom level.
	for l := s.level - 1; l >= 0; l-- {
		// Iterate value util value's index is >= given index.
		// The max iterate count is skip list's length. So the worst O(n) is N.
		for currentNode.nextNodes[l] != s.tail && currentNode.nextNodes[l].index < index {
			currentNode = currentNode.nextNodes[l]
		}

		// When next value's index is >= given index, add current value whose index < given index.
		previousNodes[l] = currentNode
	}

	// Avoid point to tail which will occur panic in Insert and Delete function.
	// When the next value is tail.
	// The index is larger than the maximum index in the skip list or skip list's length is 0. Don't point to tail.
	// When the next value isn't tail.
	// Next value's index must >= given index. Point to it.
	if currentNode.nextNodes[0] != s.tail {
		currentNode = currentNode.nextNodes[0]
	}
	// fmt.Printf("previous value:\n")
	// for _, n := range previousNodes {
	// 	fmt.Printf("%p\t", n)
	// }
	// fmt.Println()
	// fmt.Printf("end doSearch %v\n", index)

	return previousNodes, currentNode
}

// searchWithoutPreviousNodes will return the value whose index is given index.
// If can not find the given index, return nil.
// This function is faster than searchWithPreviousNodes and it used to only searching index.
func (s *skipList) searchWithoutPreviousNodes(index uint64) *Node {
	currentNode := s.head

	// Read lock and unlock.
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Iterate from top level to bottom level.
	for l := s.level - 1; l >= 0; l-- {
		// Iterate value util value's index is >= given index.
		// The max iterate count is skip list's length. So the worst O(n) is N.
		for currentNode.nextNodes[l] != s.tail && currentNode.nextNodes[l].index < index {
			currentNode = currentNode.nextNodes[l]
		}
	}

	currentNode = currentNode.nextNodes[0]
	if currentNode == s.tail || currentNode.index > index {
		return nil
	} else if currentNode.index == index {
		return currentNode
	} else {
		return nil
	}
}

// insert will insert a value into skip list and update the length.
// If skip has these this index, overwrite the value, otherwise add it.
func (s *skipList) insert(index uint64, value interface{}) {
	previousNodes, currentNode := s.searchWithPreviousNodes(index)

	// Write lock and unlock.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if currentNode != s.head && currentNode.index == index {
		currentNode.value = value
		return
	}

	// Make a new value.
	newNode := newNode(index, value, s.randomLevel())

	// Adjust pointer. Similar to update linked list.
	for i := len(newNode.nextNodes) - 1; i >= 0; i-- {
		// Firstly, new value point to next value.
		newNode.nextNodes[i] = previousNodes[i].nextNodes[i]

		// Secondly, previous nodes point to new value.
		previousNodes[i].nextNodes[i] = newNode
	}

	atomic.AddInt32(&s.length, 1)
}

// delete will find the index is existed or not firstly.
// If existed, delete it and update length, otherwise do nothing.
func (s *skipList) delete(index uint64) {
	previousNodes, currentNode := s.searchWithPreviousNodes(index)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// If skip list length is 0 or could not find value with the given index.
	if currentNode != s.head && currentNode.index == index {
		// Adjust pointer. Similar to update linked list.
		for i := 0; i < len(currentNode.nextNodes); i++ {
			previousNodes[i].nextNodes[i] = currentNode.nextNodes[i]
			currentNode.nextNodes[i] = nil
		}

		atomic.AddInt32(&s.length, -1)
	}
}

// snapshot will create a snapshot of the skip list and return a slice of the nodes.
func (s *skipList) snapshot() []*Node {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]*Node, s.length)
	i := 0

	currentNode := s.head.nextNodes[0]
	for currentNode != s.tail {
		node := &Node{
			index:     currentNode.index,
			value:     currentNode.value,
			nextNodes: nil,
		}

		result[i] = node
		currentNode = currentNode.nextNodes[0]
		i++
	}

	return result
}

// getLength will return the length of skip list.
func (s *skipList) getLength() int32 {
	return atomic.LoadInt32(&s.length)
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

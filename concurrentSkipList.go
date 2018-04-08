package ConcurrentSkipList

import (
	"math"
	"sync/atomic"

	"github.com/OneOfOne/xxhash"
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
		for j := SHARDS - 1; j > i; j-- {
			t -= step
		}

		shardIndexes[i] = t
	}
}

type ConcurrentSkipList struct {
	skipLists []*skipList
	length    int32
	level     int
}

func NewConcurrentSkipList(level int) *ConcurrentSkipList {
	if level <= 0 || level > MAX_LEVEL {
		level = MAX_LEVEL
	}

	skipLists := make([]*skipList, SHARDS)
	for i := 0; i < SHARDS; i++ {
		skipLists[i] = newSkipList(level)
	}

	return &ConcurrentSkipList{
		skipLists: skipLists,
		length:    0,
		level:     level,
	}
}

// Level will return the level of skip list.
func (s *ConcurrentSkipList) Level() int {
	return s.level
}

// Length will return the length of skip list.
func (s *ConcurrentSkipList) Length() int32 {
	var length int32 = 0
	for _, sl := range s.skipLists {
		length += sl.getLength()
	}

	return length
}

// Search will search the skip list with the given index.
// If the index exists, return the node and true, otherwise return nil and false.
func (s *ConcurrentSkipList) Search(index uint64) (*Node, bool) {
	sl := s.skipLists[getShardIndex(index)]
	if atomic.LoadInt32(&sl.length) == 0 {
		return nil, false
	}

	result := sl.searchWithoutPreviousNodes(index)
	return result, result != nil
}

// Insert will insert a node into skip list. If skip has these this index, overwrite the value, otherwise add it.
func (s *ConcurrentSkipList) Insert(index uint64, value interface{}) {
	// Ignore nil value.
	if value == nil {
		return
	}

	sl := s.skipLists[getShardIndex(index)]
	sl.insert(index, value)
}

func (s *ConcurrentSkipList) Delete(index uint64) {
	sl := s.skipLists[getShardIndex(index)]
	if atomic.LoadInt32(&sl.length) == 0 {
		return
	}

	sl.delete(index)
}

// Locate which shard the given index belong to.
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

// Hash will calculate the input's hash value using xxHash algorithm.
// It can be used to calculate the index of skip list.
// See more detail in https://cyan4973.github.io/xxHash/
func Hash(input []byte) uint64 {
	h := xxhash.New64()
	h.Write(input)
	return h.Sum64()
}

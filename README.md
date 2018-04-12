# ConcurrentSkipList
[![Build Status](https://travis-ci.org/AceDarkknight/ConcurrentSkipList.svg?branch=master)](https://travis-ci.org/AceDarkknight/ConcurrentSkipList)

A light, high-performance, concurrent, thread-safe skip list implementation written in Golang.

## Getting Start
- **Import package**

```bash
go get github.com/AceDarkknight/ConcurrentSkipList
```
```go
import "github.com/AceDarkknight/ConcurrentSkipList"
```
And you can use the package now!

- **Usage**
```go
// Create a new skip list. The parameter is the level of the skip list.
// Parameter must > 0 and <=32, if not, use the default value:32.
skipList := ConcurrentSkipList.NewConcurrentSkipList(12)

// Insert index and value. The index must uint64 and value is interface.
skipList.Insert(uint64(1), 1)
skipList.Insert(uint64(2), 2)

// Search in skip list.
if node, ok := skipList.Search(uint64(1)); ok {
	fmt.Printf("index:%v value:%v\n", node.Index(), node.Value())
}

// Delete by index.
skipList.Delete(uint64(2))

// Get the level of skip list.
_ = skipList.Level()

// Get the length of skip list.
_ = skipList.Length()

// Iterate each node in skip list.
skipList.ForEach(func(node *ConcurrentSkipList.Node) bool {
	fmt.Printf("index:%v value:%v\n", node.Index(), node.Value())
	return true
})
```

## References
https://en.wikipedia.org/wiki/Skip_list

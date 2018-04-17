package ConcurrentSkipList

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewConcurrentSkipList(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{-1}},
		{"test2", args{64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := NewConcurrentSkipList(tt.args.level); got != nil || err == nil {
				t.Errorf("NewConcurrentSkipList() = %#v,%#v", got, err)
			}
		})
	}
}

func TestConcurrentSkipList_Search(t *testing.T) {
	concurrentSkipList1, _ := NewConcurrentSkipList(16)

	type args struct {
		input uint64
	}
	type want struct {
		existed bool
		value   interface{}
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{"test1", args{uint64(math.MaxUint64)}, want{false, nil}},
		{"test2", args{0}, want{false, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, existed := concurrentSkipList1.Search(tt.args.input); existed != tt.want.existed || got != nil {
				t.Errorf("Search() = value:%v existed:%v, want value:%v existed:%v", got, existed, tt.want.value, tt.want.existed)
			}
		})
	}

	// math.MaxUint64 is untyped constant, need to convert type first.
	// See more in https://stackoverflow.com/questions/16474594/how-can-i-print-out-an-constant-uint64-in-go-using-fmt
	concurrentSkipList1.Insert(uint64(math.MaxUint64), uint64(math.MaxUint64))
	concurrentSkipList1.Insert(0, 0)

	tests = []struct {
		name string
		args args
		want want
	}{
		{"test3", args{uint64(math.MaxUint64)}, want{true, uint64(math.MaxUint64)}},
		{"test4", args{0}, want{true, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, existed := concurrentSkipList1.Search(tt.args.input); existed != tt.want.existed || got.value != tt.want.value {
				t.Errorf("Search() = value:%v existed:%v, want value:%v existed:%v", got.Value(), existed, tt.want.value, tt.want.existed)
			}
		})
	}

	concurrentSkipList1.Insert(1, nil)
	tests = []struct {
		name string
		args args
		want want
	}{
		{"test5", args{1}, want{false, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, existed := concurrentSkipList1.Search(tt.args.input); existed != tt.want.existed || got != nil {
				t.Errorf("Search() = value value:%v existed:%v, want value:%v existed:%v", got.Value(), existed, tt.want.value, tt.want.existed)
			}
		})
	}

	concurrentSkipList2, _ := NewConcurrentSkipList(8)
	for i, v := range shardIndexes {
		concurrentSkipList2.Insert(v, i)
	}

	for i, v := range shardIndexes {
		t.Run(fmt.Sprintf("test%d", i+6), func(t *testing.T) {
			if got, existed := concurrentSkipList2.Search(v); !existed || got.Value() != i {
				t.Errorf("Search() = value:%v existed:%v, want value:%v", got.Value(), existed, i)
			}
		})
	}
}

func TestConcurrentSkipList_Insert(t *testing.T) {
	skipList, _ := NewConcurrentSkipList(8)
	type args struct {
		index uint64
	}

	t.Run("test level", func(t *testing.T) {
		if skipList.Level() != 8 {
			t.Errorf("skip list level is not correct")
		}
	})

	t.Run("test Length0", func(t *testing.T) {
		if skipList.Length() != 0 {
			t.Errorf("skip list error length are not correct")
		}
	})

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"test0", args{1}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := skipList.Search(tt.args.index); ok || got != nil {
				t.Errorf("Search() = %v, want = %v", got, tt.want)
			}
		})
	}

	// Insert elements.
	for i := 2; i <= 10; i++ {
		skipList.Insert(uint64(i), i)
	}

	skipList.Insert(uint64(1), 1)
	skipList.Insert(uint64(2018), 111)

	t.Run("test Length1", func(t *testing.T) {
		if length := skipList.Length(); length != 11 {
			t.Errorf("skip list's length is not correct, got %d", length)
		}
	})

	tests = []struct {
		name string
		args args
		want interface{}
	}{
		{"test1", args{1}, 1},
		{"test2", args{2}, 2},
		{"test3", args{3}, 3},
		{"test4", args{4}, 4},
		{"test5", args{5}, 5},
		{"test6", args{6}, 6},
		{"test7", args{7}, 7},
		{"test8", args{8}, 8},
		{"test9", args{9}, 9},
		{"test10", args{10}, 10},
		{"test11", args{2018}, 111},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := skipList.Search(tt.args.index); !ok || got.Value() != tt.want {
				t.Errorf("Search() = %v, want = %v", got, tt.want)
			}
		})
	}

	skipList.Insert(uint64(1024), nil)
	skipList.Insert(uint64(2018), 2017)
	t.Run("test length2", func(t *testing.T) {
		if got := skipList.Length(); got != 11 {
			t.Errorf("skip list's length = %d is not correct", got)
		}
	})

	tests = []struct {
		name string
		args args
		want interface{}
	}{
		{"test13", args{1}, 1},
		{"test14", args{2}, 2},
		{"test15", args{3}, 3},
		{"test16", args{4}, 4},
		{"test17", args{5}, 5},
		{"test18", args{6}, 6},
		{"test19", args{7}, 7},
		{"test20", args{8}, 8},
		{"test21", args{9}, 9},
		{"test22", args{10}, 10},
		{"test23", args{2018}, 2017},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := skipList.Search(tt.args.index); !ok || got.Value() != tt.want {
				t.Errorf("Search() = %v, want = %v", got, tt.want)
			}
		})
	}
	t.Run("test24", func(t *testing.T) {
		if got, ok := skipList.Search(1024); ok && got != nil {
			t.Errorf("Search() = %v, want = %v", got, nil)
		}
	})

	skipList.Insert(uint64(1<<60), uint64(1<<60))
	skipList.Insert(uint64(1<<61), uint64(1<<61))
	skipList.Insert(uint64(1<<62), uint64(1<<62))

	var lastIndex uint64
	skipList.ForEach(func(node *Node) bool {
		t.Run("test sequence", func(t *testing.T) {
			t.Logf("index:%v value:%v", node.Index(), node.Value())
			if lastIndex > node.Index() {
				t.Errorf("incorrect sequence")
			}

			lastIndex = node.Index()
		})

		return true
	})
}

func TestConcurrentSkipList_Delete(t *testing.T) {
	skipList, _ := NewConcurrentSkipList(16)
	skipList.Delete(uint64(1))
	type args struct {
		index uint64
	}

	t.Run("test level", func(t *testing.T) {
		if skipList.Level() != 16 {
			t.Errorf("skip list's level is not correct")
		}
	})

	t.Run("test Length0", func(t *testing.T) {
		if length := skipList.Length(); length != 0 {
			t.Errorf("skip list's length is not correct, got %d", length)
		}
	})

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"test0", args{0}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := skipList.Search(tt.args.index); ok || got != nil {
				t.Errorf("Search() = %v, want = %v", got, tt.want)
			}
		})
	}

	// Insert elements.
	for i := 0; i <= 10; i++ {
		skipList.Insert(uint64(i), i)
	}

	t.Run("test Length1", func(t *testing.T) {
		if length := skipList.Length(); length != 11 {
			t.Errorf("skip list's length is not correct, got %d", length)
		}
	})

	// Delete elements.
	skipList.Delete(uint64(5))
	skipList.Delete(uint64(1))
	skipList.Delete(uint64(11))

	t.Run("test Length1", func(t *testing.T) {
		if length := skipList.Length(); length != 9 {
			t.Errorf("skip list's length is not correct, got %d", length)
		}
	})

	tests = []struct {
		name string
		args args
		want interface{}
	}{
		{"test1", args{2}, 2},
		{"test2", args{3}, 3},
		{"test3", args{4}, 4},
		{"test4", args{6}, 6},
		{"test5", args{7}, 7},
		{"test6", args{8}, 8},
		{"test7", args{9}, 9},
		{"test8", args{10}, 10},
		{"test9", args{0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := skipList.Search(tt.args.index); !ok || got.Value() != tt.want {
				t.Errorf("Search() = %v, want = %v", got, tt.want)
			}
		})
	}

	tests = []struct {
		name string
		args args
		want interface{}
	}{
		{"test10", args{1}, nil},
		{"test11", args{5}, nil},
		{"test12", args{11}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := skipList.Search(tt.args.index); ok || got != nil {
				t.Errorf("Search() = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentSkipList_ForEach(t *testing.T) {
	skipList, _ := NewConcurrentSkipList(10)
	indexes := make([]uint64, 0)
	count := 10000
	for i := 0; i < count; i++ {
		index := Hash([]byte(strconv.Itoa(i)))
		indexes = append(indexes, index)
	}

	for i, index := range indexes {
		skipList.Insert(index, i)
	}

	i := -1
	skipList.ForEach(func(node *Node) bool {
		i++
		return false
	})

	t.Run("test", func(t *testing.T) {
		if i != 0 {
			t.Errorf("ForEach() occur error. got %d", i)
		}
	})
}

func TestConcurrentSkipList_Insert_Parallel(t *testing.T) {
	skipList, _ := NewConcurrentSkipList(10)
	indexes := make([]uint64, 0)
	count := 10000
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		index := Hash([]byte(strconv.Itoa(i)))
		indexes = append(indexes, index)
	}

	for _, index := range indexes {
		wg.Add(1)
		go func(i uint64, v interface{}) {
			skipList.Insert(i, v)
			wg.Done()
		}(index, index)
	}

	wg.Wait()
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})

	t.Run("test length", func(t *testing.T) {
		if length := skipList.Length(); length != int32(count) {
			t.Errorf("skip list's length are not correct, got %d", length)
		}
	})

	for i, index := range indexes {
		t.Run(fmt.Sprintf("test %d", i+1), func(t *testing.T) {
			if got, ok := skipList.Search(index); !ok || got.Index() != index || got.Value() != index {
				t.Errorf("Search() = %v, want = %v", got, index)
			}
		})
	}
}

func TestConcurrentSkipList_Delete_Parallel(t *testing.T) {
	skipList, _ := NewConcurrentSkipList(10)
	indexes := make([]uint64, 0)
	count := 10000
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		index := Hash([]byte(strconv.Itoa(i)))
		indexes = append(indexes, index)
	}

	for _, index := range indexes {
		wg.Add(1)
		go func(i uint64, v interface{}) {
			skipList.Insert(i, v)
			wg.Done()
		}(index, index)
	}

	wg.Wait()
	t.Run("test length1", func(t *testing.T) {
		if length := skipList.Length(); length != int32(count) {
			t.Errorf("skip list's length are not correct, got %d", length)
		}
	})

	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})
	for _, index := range indexes {
		wg.Add(1)
		go func(i uint64, v interface{}) {
			skipList.Delete(i)
			wg.Done()
		}(index, index)
	}
	for _, index := range indexes {
		wg.Add(1)
		go func(i uint64, v interface{}) {
			skipList.Delete(i)
			wg.Done()
		}(index, index)
	}
	wg.Wait()

	t.Run("test length2", func(t *testing.T) {
		if length := skipList.Length(); length != 0 {
			t.Errorf("skip list's length are not correct, got %d", length)
		}
	})
}

func TestConcurrentSkipList_ForEach_Parallel(t *testing.T) {
	skipList, _ := NewConcurrentSkipList(10)
	indexes := make([]uint64, 0)
	count := 10000
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		indexes = append(indexes, uint64(i))
	}

	for _, index := range indexes {
		wg.Add(1)
		go func(i uint64, v interface{}) {
			skipList.Insert(i, v)
			wg.Done()
		}(index, index)
	}

	wg.Wait()
	go func() {
		i := 0
		skipList.ForEach(func(node *Node) bool {
			i++
			return true
		})

		t.Run("test1", func(t *testing.T) {
			if i != count {
				t.Errorf("skip list's ForEach() occurs error, node count got %d,want %d", i, count)
			}
		})
	}()

	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})

	time.Sleep(time.Second)
	for _, index := range indexes {
		wg.Add(1)
		go func(i uint64, v interface{}) {
			skipList.Delete(i)
			wg.Done()
		}(index, index)
	}
	wg.Wait()

	i := 0
	skipList.ForEach(func(node *Node) bool {
		i++
		return true
	})

	t.Run("test1", func(t *testing.T) {
		if i != 0 {
			t.Errorf("skip list's ForEach() occurs error, node count got %d,want %d", i, count)
		}
	})
}

func TestConcurrentSkipList_Level(t *testing.T) {
	skipList, _ := NewConcurrentSkipList(16)
	for i := 0; i < 100000; i++ {
		index := Hash([]byte(strconv.Itoa(i)))
		skipList.Insert(index, i)
	}

	length := skipList.Length()
	levels := make([]int, 17)
	for _, sl := range skipList.skipLists {
		if sl.getLength() == 0 {
			continue
		}

		currentNode := sl.head
		for currentNode.nextNodes[0] != sl.tail {
			levels[len(currentNode.nextNodes)]++
			currentNode = currentNode.nextNodes[0]
		}
	}

	fmt.Printf("level count:%#v\n", levels)
	t.Run("test", func(t *testing.T) {
		got := 0
		for _, l := range levels {
			got += l
		}

		if got != int(length) {
			t.Errorf("count of each level= %v, want %v", got, length)
		}
	})
}

func TestConcurrentSkipList_Sub(t *testing.T) {
	skipList, _ := NewConcurrentSkipList(12)
	count := 100
	indexes := make([]uint64, count)
	for i := 0; i < count; i++ {
		index := Hash([]byte(strconv.Itoa(i)))
		indexes[i] = index
		skipList.Insert(index, index)
	}

	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})

	type args struct {
		startIndex int32
		length     int32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"test1", args{0, -1}, 0},
		{"test2", args{-1, -1}, 0},
		{"test3", args{105, 2}, 0},
		{"test4", args{100, 2}, 0},
		{"test5", args{99, 55}, 1},
		{"test6", args{51, 55}, 49},
		{"test7", args{0, 100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := skipList.Sub(tt.args.startIndex, tt.args.length); len(got) != tt.want {
				t.Errorf("Sub() = %v, want = %v", len(got), tt.want)
			}
		})
	}

	got := skipList.Sub(0, 100)
	for i := 0; i < len(indexes); i++ {
		t.Run(fmt.Sprintf("test%d", i+8), func(t *testing.T) {
			if got[i].Index() != indexes[i] {
				t.Errorf("Sub() = %v, want = %v", got[i].Index(), indexes[i])
			}
		})
	}
}

func TestHash(t *testing.T) {
	input := `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
`

	type args struct {
		input []byte
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"test", args{[]byte(input)}, 0xFFAE31BEBFED7652},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.input); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

package ConcurrentSkipList

import (
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkSkipList_Insert_Ordered(b *testing.B) {
	skipList := NewConcurrentSkipList(12)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.Intn(b.N)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Insert(uint64(i), index)
	}
}

func BenchmarkSkipList_Insert_Randomly(b *testing.B) {
	skipList := NewConcurrentSkipList(12)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.Intn(b.N)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Insert(index, index)
	}
}

func BenchmarkSkipList_Search_100000Elements(b *testing.B) {
	skipList := NewConcurrentSkipList(32)
	for i := 0; i < 100000; i++ {
		skipList.Insert(uint64(i), i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.Intn(b.N)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Search(index)
	}
}

func BenchmarkSkipList_Search_1000000Elements(b *testing.B) {
	skipList := NewConcurrentSkipList(32)
	for i := 0; i < 1000000; i++ {
		skipList.Insert(uint64(i), i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.Intn(b.N)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Search(index)
	}
}

func BenchmarkSkipList_Search_10000000Elements(b *testing.B) {
	skipList := NewConcurrentSkipList(32)
	for i := 0; i < 10000000; i++ {
		skipList.Insert(uint64(i), i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.Intn(b.N)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Search(index)
	}
}

func BenchmarkSkipList_Search_12Level(b *testing.B) {
	skipList := NewConcurrentSkipList(12)
	for i := 0; i < 10000000; i++ {
		skipList.Insert(uint64(i), i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.Intn(b.N)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Search(index)
	}
}

func BenchmarkSkipList_Search_24Level(b *testing.B) {
	skipList := NewConcurrentSkipList(24)
	for i := 0; i < 10000000; i++ {
		skipList.Insert(uint64(i), i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.Intn(b.N)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Search(index)
	}
}

func BenchmarkConcurrentSkipList_Insert_Parallel(b *testing.B) {
	skipList := NewConcurrentSkipList(12)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			skipList.Insert(uint64(rand.Intn(b.N)), 0)
		}
	})
}

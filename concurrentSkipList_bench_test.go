package ConcurrentSkipList

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkSkipList_Insert_Ordered(b *testing.B) {
	skipList := NewConcurrentSkipList(10)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000)
		_ = Hash([]byte(strconv.Itoa(t)))
		skipList.Insert(uint64(i), i)
	}
}

func BenchmarkSkipList_Insert_Randomly(b *testing.B) {
	skipList := NewConcurrentSkipList(10)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100000)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Insert(index, i)
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
		t := rand.Intn(100000)
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
		t := rand.Intn(1000000)
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
		t := rand.Intn(10000000)
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
		t := rand.Intn(10000000)
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
		t := rand.Intn(10000000)
		index := Hash([]byte(strconv.Itoa(t)))
		skipList.Search(index)
	}
}

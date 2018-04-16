package ConcurrentSkipList

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkConcurrentSkipList_Insert_Ordered(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(12)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Hash([]byte(strconv.Itoa(i)))
		skipList.Insert(uint64(i), i)
	}
}

func BenchmarkConcurrentSkipList_Insert_Randomly(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(12)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		index := Hash([]byte(strconv.Itoa(i)))
		skipList.Insert(index, index)
	}
}
func BenchmarkConcurrentSkipList_Delete(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(12)
	for i := 0; i < 10000000; i++ {
		skipList.Insert(uint64(i), i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// t := rand.Intn(b.N)
		// index := Hash([]byte(strconv.Itoa(t)))
		skipList.Delete(uint64(i))
	}
}

func BenchmarkConcurrentSkipList_Search_100000Elements(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(12)
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

func BenchmarkConcurrentSkipList_Search_1000000Elements(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(12)
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

func BenchmarkConcurrentSkipList_Search_10000000Elements(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(12)
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

func BenchmarkConcurrentSkipList_Search_24Level(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(24)
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

func BenchmarkConcurrentSkipList_Search_32Level(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(32)
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
	skipList, _ := NewConcurrentSkipList(12)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			skipList.Insert(uint64(rand.Intn(b.N)), 0)
		}
	})
}

func BenchmarkConcurrentSkipList_Delete_Parallel(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(12)
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < b.N; i++ {
			t := rand.Intn(i + 1)
			index := Hash([]byte(strconv.Itoa(t)))
			wg.Add(1)
			go func() {
				defer wg.Done()
				skipList.Insert(index, index)
			}()
		}

		wg.Wait()
	}()

	go func() {
		var wg sync.WaitGroup
		for i := 0; i < b.N; i++ {
			t := rand.Intn(b.N)
			index := Hash([]byte(strconv.Itoa(t)))
			wg.Add(1)
			go func() {
				defer wg.Done()
				skipList.Search(index)
			}()
		}

		wg.Wait()
	}()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			skipList.Delete(uint64(rand.Intn(b.N)))
		}
	})
}

func BenchmarkConcurrentSkipList_Search_Parallel(b *testing.B) {
	skipList, _ := NewConcurrentSkipList(12)
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < b.N; i++ {
			t := rand.Intn(b.N)
			index := Hash([]byte(strconv.Itoa(t)))
			wg.Add(1)
			go func() {
				defer wg.Done()
				skipList.Insert(index, index)
			}()
		}

		wg.Wait()
	}()

	go func() {
		var wg sync.WaitGroup
		for i := 0; i < b.N; i++ {
			t := rand.Intn(b.N)
			index := Hash([]byte(strconv.Itoa(t)))
			wg.Add(1)
			go func() {
				defer wg.Done()
				skipList.Delete(index)
			}()
		}

		wg.Wait()
	}()

	time.Sleep(time.Millisecond)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			skipList.Search(uint64(rand.Intn(b.N)))
		}
	})
}

func BenchmarkHash(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Hash([]byte(strconv.Itoa(i)))
	}
}

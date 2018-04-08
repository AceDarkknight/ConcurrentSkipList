package ConcurrentSkipList

import (
	"math"
	"testing"
)

func TestConcurrentSkipList_Search_SingleThread(t *testing.T) {
	skipList := NewConcurrentSkipList(16)

	type args struct {
		input uint64
	}
	tests1 := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{math.MaxUint64}, false},
		{"test2", args{0}, false},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := skipList.Search(tt.args.input); got != tt.want {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}

	skipList.Insert(math.MaxUint64, math.MaxUint64)
	skipList.Insert(0, 0)
}

func TestConcurrentSkipList_Insert_SingleThread(t *testing.T) {

}

package ConcurrentSkipList

import (
	"math"
	"testing"
)

func TestConcurrentSkipList_Search_SingleThread(t *testing.T) {
	concurrentSkipList := NewConcurrentSkipList(16)

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
			if got, existed := concurrentSkipList.Search(tt.args.input); existed != tt.want.existed || got != tt.want.value {
				t.Errorf("Search() = node value:%v existed:%v, want value:%v existed:%v", got, existed, tt.want.value, tt.want.existed)
			}
		})
	}

	// math.MaxUint64 is constant, need to convert type first.
	// See more in https://stackoverflow.com/questions/16474594/how-can-i-print-out-an-constant-uint64-in-go-using-fmt
	concurrentSkipList.Insert(uint64(math.MaxUint64), uint64(math.MaxUint64))
	concurrentSkipList.Insert(0, 0)

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
			if got, existed := concurrentSkipList.Search(tt.args.input); existed != tt.want.existed || got.value != tt.want.value {
				t.Errorf("Search() = node value:%v existed:%v, want value:%v existed:%v", got.Value(), existed, tt.want.value, tt.want.existed)
			}
		})
	}
}

func TestConcurrentSkipList_Insert_SingleThread(t *testing.T) {

}

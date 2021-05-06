package lru

import (
	"strconv"
	"testing"
)

func TestLRU_Capacity(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: 127, want: 128},
		{input: 128, want: 128},
		{input: 129, want: 256},
	}

	for _, tc := range tests {
		l, err := New(WithCapacity(tc.input))
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		got := l.Cap()
		if got != tc.want {
			t.Fatalf("Cap expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestLRU_Concurrency(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: 15, want: 16},
		{input: 16, want: 16},
		{input: 17, want: 32},
	}

	for _, tc := range tests {
		l, err := New(WithConcurrency(tc.input))
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		got := len(l.(*cache).segments)
		if got != tc.want {
			t.Fatalf("Concurrency expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestLRU_SegmentBalance(t *testing.T) {
	l, err := New(WithCapacity(8192))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i := int64(0); i < 8192; i++ {
		l.Set(strconv.FormatInt(i, 10), i)
	}
	//segment balance
	segments := l.(*cache).segments
	threshold := 0.5
	for _, each := range segments {
		maxlen := float64(l.Len()) * (1.0 + float64(threshold)) / float64(len(segments))
		minlen := float64(l.Len()) * (1.0 - float64(threshold)) / float64(len(segments))
		got := float64(each.Len())
		if got < minlen || got > maxlen {
			t.Fatalf("segment len expected: %v, %v , got: %v", minlen, maxlen, got)
		}
	}
}

func TestLRU_Add(t *testing.T) {
	l, err := New(WithCapacity(8192), WithConcurrency(1))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i := int64(0); i < 2*8192; i++ {
		l.Set(strconv.FormatInt(i, 10), i)
	}
	threshold := 0.5
	for i := int64(0); i < int64(8192*(1-threshold)); i++ {
		if _, ok := l.Get(strconv.FormatInt(i, 10)); ok {
			t.Fatalf("should not exist:%v", i)
		}
	}
	for i := int64(8192 * (1 + threshold)); i < 2*8192; i++ {
		if _, ok := l.Get(strconv.FormatInt(i, 10)); !ok {
			t.Fatalf("should exist:%v", i)
		}
	}
}

func TestLRU_Delete(t *testing.T) {
	l, err := New(WithCapacity(8192), WithConcurrency(1))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i := int64(0); i < 2*8192; i++ {
		l.Set(strconv.FormatInt(i, 10), i)
	}
	threshold := 0.5
	for i := int64(8192 * (1 + threshold)); i < 2*8192; i++ {
		l.Delete(strconv.FormatInt(i, 10))
	}
	for i := int64(8192 * (1 + threshold)); i < 2*8192; i++ {
		if _, ok := l.Get(strconv.FormatInt(i, 10)); ok {
			t.Fatalf("should not exist:%v", i)
		}
	}
}

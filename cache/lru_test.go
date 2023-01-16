package cache

import (
	"strconv"
	"testing"
)

func TestLRU_Capacity(t *testing.T) {
	tests := []struct {
		input   int
		want    int
		wantErr bool
	}{
		{input: 0, want: 0, wantErr: true},
		{input: 127, want: 128, wantErr: false},
		{input: 128, want: 128, wantErr: false},
		{input: 129, want: 256, wantErr: false},
		{input: 1<<30 + 1, want: 1 << 30, wantErr: false},
	}

	for _, tc := range tests {
		l, err := NewLRU(WithCapacity(tc.input))
		gotErr := (err != nil)
		if gotErr != tc.wantErr {
			t.Fatalf("Cap expected err: %v, got:%v", tc.wantErr, err)
		}

		if err != nil {
			continue
		}

		got := l.Cap()
		if got != tc.want {
			t.Fatalf("Cap expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestLRU_Concurrency(t *testing.T) {
	tests := []struct {
		input   int
		want    int
		wantErr bool
	}{
		{input: 0, want: 0, wantErr: true},
		{input: 15, want: 16, wantErr: false},
		{input: 16, want: 16, wantErr: false},
		{input: 17, want: 32, wantErr: false},
		{input: 65537, want: 65536, wantErr: false},
	}

	for _, tc := range tests {
		l, err := NewLRU(WithConcurrency(tc.input))
		gotErr := (err != nil)
		if gotErr != tc.wantErr {
			t.Fatalf("Concurrency expected err: %v, got:%v", tc.wantErr, err)
		}

		if err != nil {
			continue
		}

		got := len(l.segments)
		if got != tc.want {
			t.Fatalf("Concurrency expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestLRU_SegmentBalance(t *testing.T) {
	l, err := NewLRU(WithCapacity(8192))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i := int64(0); i < 8192; i++ {
		l.Set(strconv.FormatInt(i, 10), i)
	}
	// segment balance
	segments := l.segments
	threshold := 0.5
	for _, each := range segments {
		maxlen := float64(l.Len()) * (1.0 + threshold) / float64(len(segments))
		minlen := float64(l.Len()) * (1.0 - threshold) / float64(len(segments))
		got := float64(each.Len())
		if got < minlen || got > maxlen {
			t.Fatalf("segment len expected: %v, %v , got: %v", minlen, maxlen, got)
		}
	}
}

func TestLRU_SetGet(t *testing.T) {
	l, err := NewLRU(WithCapacity(8192), WithConcurrency(1))
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
	l, err := NewLRU(WithCapacity(8192), WithConcurrency(1))
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

func TestLRU_Exist(t *testing.T) {
	l, err := NewLRU(WithCapacity(8192), WithConcurrency(1))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i := int64(0); i < 2*8192; i++ {
		l.Set(strconv.FormatInt(i, 10), i)
	}
	threshold := 0.5
	for i := int64(0); i < int64(8192*(1-threshold)); i++ {
		if ok := l.Exists(strconv.FormatInt(i, 10)); ok {
			t.Fatalf("should not exist:%v", i)
		}
	}
	for i := int64(8192 * (1 + threshold)); i < 2*8192; i++ {
		if ok := l.Exists(strconv.FormatInt(i, 10)); !ok {
			t.Fatalf("should exist:%v", i)
		}
	}
}

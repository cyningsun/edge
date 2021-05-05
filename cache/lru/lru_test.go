package lru

import (
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

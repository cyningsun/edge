package cache

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/cyningsun/edge/internal/cache/lru"
)

type bitwise struct {
	concurrency int
	capacity    int

	normal *normalize
}

var testSizeData = map[string]*bitwise{
	"normal": {
		16,
		8192,
		&normalize{
			16,
			512,
			28,
			15,
		},
	},
	"concurrency": {
		15,
		8192,
		&normalize{
			16,
			512,
			28,
			15,
		},
	},
	"capacity": {
		16,
		8000,
		&normalize{
			16,
			512,
			28,
			15,
		},
	},
	"max": {
		maxSegments,
		maxCapacity,
		&normalize{
			65536,
			16384,
			16,
			65535,
		},
	},
	"maxConcurrency": {
		maxSegments + 1,
		maxCapacity,
		&normalize{
			65536,
			16384,
			16,
			65535,
		},
	},
	"maxCapacity": {
		maxSegments,
		maxCapacity + 1,
		&normalize{
			65536,
			16384,
			16,
			65535,
		},
	},
}

func Test_bitwiseOpt(t *testing.T) {
	type args struct {
		concurrency int
		capacity    int
	}
	tests := []struct {
		name  string
		args  args
		size  int
		cap   int
		shift uint32
		mask  uint32
	}{
		{
			"normal",
			args{
				testSizeData["normal"].concurrency,
				testSizeData["normal"].capacity,
			},
			testSizeData["normal"].normal.size,
			testSizeData["normal"].normal.cap,
			testSizeData["normal"].normal.shift,
			testSizeData["normal"].normal.mask,
		},
		{
			"max",
			args{
				testSizeData["max"].concurrency,
				testSizeData["max"].capacity,
			},
			testSizeData["max"].normal.size,
			testSizeData["max"].normal.cap,
			testSizeData["max"].normal.shift,
			testSizeData["max"].normal.mask,
		},
		{
			"normalize concurrency",
			args{
				testSizeData["concurrency"].concurrency,
				testSizeData["concurrency"].capacity,
			},
			testSizeData["concurrency"].normal.size,
			testSizeData["concurrency"].normal.cap,
			testSizeData["concurrency"].normal.shift,
			testSizeData["concurrency"].normal.mask,
		},
		{
			"normalize capacity",
			args{
				testSizeData["capacity"].concurrency,
				testSizeData["capacity"].capacity,
			},
			testSizeData["capacity"].normal.size,
			testSizeData["capacity"].normal.cap,
			testSizeData["capacity"].normal.shift,
			testSizeData["capacity"].normal.mask,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normal := bitwiseOpt(tt.args.concurrency, tt.args.capacity)
			if normal.size != tt.size {
				t.Errorf("bitwiseOpt() size = %v, want %v", normal.size, tt.size)
			}
			if normal.cap != tt.cap {
				t.Errorf("bitwiseOpt() cap = %v, want %v", normal.cap, tt.cap)
			}
			if normal.shift != tt.shift {
				t.Errorf("bitwiseOpt() shift = %v, want %v", normal.shift, tt.shift)
			}
			if normal.mask != tt.mask {
				t.Errorf("bitwiseOpt() mask = %v, want %v", normal.mask, tt.mask)
			}
		})
	}
}

func TestNewLRU(t *testing.T) {
	type args struct {
		opts []Opt
	}
	tests := []struct {
		name      string
		args      args
		wantParam *bitwise
		wantErr   bool
	}{
		{
			"normal",
			args{},
			testSizeData["normal"],
			false,
		},
		{
			"max concurrency",
			args{
				[]Opt{
					WithConcurrency(testSizeData["maxConcurrency"].concurrency),
					WithCapacity(testSizeData["maxConcurrency"].capacity),
				},
			},
			testSizeData["maxConcurrency"],
			false,
		},
		{
			"max capacity",
			args{
				[]Opt{
					WithConcurrency(testSizeData["maxCapacity"].concurrency),
					WithCapacity(testSizeData["maxCapacity"].capacity),
				},
			},
			testSizeData["maxCapacity"],
			false,
		},
		{
			"invalid min concurrency",
			args{
				[]Opt{
					WithConcurrency(0),
				},
			},
			nil,
			true,
		},
		{
			"invalid min capacity",
			args{
				[]Opt{
					WithCapacity(0),
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var want *cache
			if tt.wantParam != nil {
				normal := tt.wantParam.normal
				segments := make([]*lru.Segment, normal.size)
				for i := range segments {
					segments[i] = lru.NewSegment(normal.cap)
				}
				want = &cache{
					segments:     segments,
					segmentMask:  normal.mask,
					segmentShift: normal.shift,
					capacity:     normal.cap * normal.size,
				}
			}

			got, err := NewLRU(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLRU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("NewLRU() = %v, want %v", got, want)
			}
		})
	}
}

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
		if ok := l.Exists(strconv.FormatInt(i, 10)); !ok {
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

// Package lru implements lru algorithm using linked list and hash map
package cache

import (
	"errors"
	"hash/fnv"

	"github.com/cyningsun/edge"
	"github.com/cyningsun/edge/internal/cache/lru"
)

var _ edge.Cache = &cache{}

// cache is concurrent safe lru cache.
// It using multi-segment to minimize RWMutex impact on performance
type cache struct {
	segments     []*lru.Segment
	segmentMask  uint32
	segmentShift uint32
	capacity     int
}

type normalize struct {
	size  int
	cap   int
	shift uint32
	mask  uint32
}

func NewLRU(opts ...Opt) (*cache, error) {
	options := &options{
		concurrency: 16,
		capacity:    8192,
	}
	for _, each := range opts {
		each(options)
	}

	switch {
	case options.capacity <= 0:
		return nil, errors.New("lru capacity invalid")
	case options.concurrency <= 0:
		return nil, errors.New("lru concurrency invalid")
	}

	if options.concurrency > maxSegments {
		options.concurrency = maxSegments
	}

	if options.capacity > maxCapacity {
		options.capacity = maxCapacity
	}

	// Find power-of-two sizes best matching arguments
	normalize := bitwiseOpt(options.concurrency, options.capacity)

	segments := make([]*lru.Segment, normalize.size)
	for i := range segments {
		segments[i] = lru.NewSegment(normalize.cap)
	}
	return &cache{
		segments:     segments,
		segmentMask:  normalize.mask,
		segmentShift: normalize.shift,
		capacity:     normalize.cap * normalize.size,
	}, nil
}

func bitwiseOpt(concurrency, capacity int) *normalize {
	shift := 0
	ssize := 1
	for ssize < concurrency {
		shift++
		ssize <<= 1
	}

	sshift := uint32(32 - shift)
	smask := uint32(ssize - 1)

	c := capacity / ssize
	if c*ssize < capacity {
		c++
	}

	scap := 1
	for scap < c {
		scap <<= 1
	}
	return &normalize{ssize, scap, sshift, smask}
}

func (c *cache) Set(key string, val interface{}) interface{} {
	seg := c.segmentFor(key)
	return seg.Set(key, val)
}

func (c *cache) Get(key string) (value interface{}, ok bool) {
	seg := c.segmentFor(key)
	return seg.Get(key)
}

func (c *cache) Delete(key string) (present bool) {
	seg := c.segmentFor(key)
	return seg.Delete(key)
}

func (c *cache) Exists(key string) bool {
	seg := c.segmentFor(key)
	return seg.Exists(key)
}

func (c *cache) Cap() int {
	return c.capacity
}

func (c *cache) Len() int {
	len := 0
	for _, each := range c.segments {
		len += each.Len()
	}
	return len
}

func (c *cache) segmentFor(key string) *lru.Segment {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	hash := h.Sum32()
	return c.segments[(hash>>c.segmentShift)&c.segmentMask]
}

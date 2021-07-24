// Package lru implements lru algorithm using linked list and hash map
package lru

import (
	"errors"
	"hash/fnv"
)

// cache is concurrent safe lru cache.
// It using multi-segment to minimize RWMutex impact on performance
type cache struct {
	segments     []*segment
	segmentMask  uint32
	segmentShift uint32
	capacity     int
}

func New(opts ...Opt) (*cache, error) {
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

	// Find power-of-two sizes best matching arguments
	sshift := 0
	ssize := 1
	for ssize < options.concurrency {
		sshift++
		ssize <<= 1
	}
	segmentShift := uint32(32 - sshift)
	segmentMask := uint32(ssize - 1)
	if options.capacity > maxCapacity {
		options.capacity = maxCapacity
	}

	c := options.capacity / ssize
	if c*ssize < options.capacity {
		c++
	}

	cap := 1
	for cap < c {
		cap <<= 1
	}

	segments := make([]*segment, ssize)
	for i := range segments {
		segments[i] = newSegment(cap)
	}
	return &cache{
		segments:     segments,
		segmentMask:  segmentMask,
		segmentShift: segmentShift,
		capacity:     cap * ssize,
	}, nil
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

func (c *cache) segmentFor(key string) *segment {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	hash := h.Sum32()
	return c.segments[(hash>>c.segmentShift)&c.segmentMask]
}

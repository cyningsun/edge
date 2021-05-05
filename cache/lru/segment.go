package lru

import (
	"container/list"
	"sync"
)

type Key interface{}

type entry struct {
	key Key
	val interface{}
}

type segment struct {
	cache map[interface{}]*list.Element
	ll    *list.List
	mtx   sync.Mutex
	cap   int
}

func newSegment(c int) *segment {
	return &segment{
		cache: make(map[interface{}]*list.Element),
		ll:    list.New(),
		cap:   c,
	}
}

func (s *segment) Add(key Key, val interface{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.cache == nil {
		s.cache = make(map[interface{}]*list.Element)
		s.ll = list.New()
	}
	if found, ok := s.cache[key]; ok {
		s.ll.MoveToFront(found)
		found.Value.(*entry).val = val
		return
	}
	new := s.ll.PushFront(&entry{key, val})
	s.cache[key] = new

	if s.cap != 0 && s.ll.Len() > s.cap {
		s.removeOldest()
	}
}

func (s *segment) Get(key Key) (val interface{}, ok bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.cache == nil {
		return
	}
	if found, hit := s.cache[key]; hit {
		s.ll.MoveToFront(found)
		return found.Value.(*entry).val, true
	}
	return
}

func (s *segment) Delete(key Key) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.cache == nil {
		return false
	}
	found, hit := s.cache[key]
	if hit {
		s.removeElement(found)
	}
	return hit
}

func (s *segment) Exists(key Key) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.cache == nil {
		return false
	}
	_, hit := s.cache[key]
	return hit
}

func (s *segment) Len() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.cache == nil {
		return 0
	}
	return s.ll.Len()
}

func (s *segment) removeOldest() {
	if s.cache == nil {
		return
	}
	found := s.ll.Back()
	if found != nil {
		s.removeElement(found)
	}
}

func (s *segment) removeElement(e *list.Element) {
	s.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(s.cache, kv.key)
}

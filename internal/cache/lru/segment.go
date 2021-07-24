package lru

import (
	"container/list"
	"sync"
)

type entry struct {
	key string
	val interface{}
}

type Segment struct {
	cache map[interface{}]*list.Element
	ll    *list.List
	mtx   sync.RWMutex
	cap   int
}

func NewSegment(c int) *Segment {
	return &Segment{
		cache: make(map[interface{}]*list.Element),
		ll:    list.New(),
		cap:   c,
	}
}

func (s *Segment) Set(key string, val interface{}) interface{} {
	m.Set.Add(1)
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.cache == nil {
		s.cache = make(map[interface{}]*list.Element)
		s.ll = list.New()
	}
	if found, ok := s.cache[key]; ok {
		s.ll.MoveToFront(found)
		oldVal := found.Value.(*entry).val
		found.Value.(*entry).val = val
		return oldVal
	}
	new := s.ll.PushFront(&entry{key, val})
	s.cache[key] = new

	if s.cap != 0 && s.ll.Len() > s.cap {
		m.Evict.Add(1)
		s.removeOldest()
	}
	return nil
}

func (s *Segment) Get(key string) (val interface{}, ok bool) {
	m.Get.Add(1)
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.cache == nil {
		return
	}
	if found, hit := s.cache[key]; hit {
		m.Hit.Add(1)
		s.ll.MoveToFront(found)
		return found.Value.(*entry).val, true
	}
	return
}

func (s *Segment) Delete(key string) bool {
	m.Delete.Add(1)
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

func (s *Segment) Exists(key string) bool {
	m.Exists.Add(1)
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.cache == nil {
		return false
	}
	_, hit := s.cache[key]
	return hit
}

func (s *Segment) Len() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.cache == nil {
		return 0
	}
	return s.ll.Len()
}

func (s *Segment) removeOldest() {
	if s.cache == nil {
		return
	}
	found := s.ll.Back()
	if found != nil {
		s.removeElement(found)
	}
}

func (s *Segment) removeElement(e *list.Element) {
	s.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(s.cache, kv.key)
}

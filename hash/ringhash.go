package hash

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
)

/*
   Consistent hashing was first described in a paper, ["Consistent hashing and random trees: Distributed caching protocols for relieving hot spots on the World Wide Web (1997)"](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.147.1879) by David Karger et al.
   It is used in distributed storage systems like Amazon Dynamo and memcached.
*/

const (
	minReplicas = 1
)

type Node interface {
	String() string
}

type ring struct {
	replicas int

	vnodes map[uint32]Node
	sorted []uint32
	mtx    sync.Mutex
}

func NewRing(replicas int) (*ring, error) {
	if replicas < minReplicas {
		return nil, fmt.Errorf("min validate replicas:%v, input:%v", minReplicas, replicas)
	}
	return &ring{
		replicas: replicas,
		vnodes:   map[uint32]Node{},
		sorted:   []uint32{},
	}, nil
}

func (r *ring) Add(node Node) {
	ringVar.Add.Add(1)

	newHash := make([]uint32, 0, r.replicas)
	for i := 1; i <= r.replicas; i++ {
		h := hash(node.String() + "_" + strconv.Itoa(i))
		newHash = append(newHash, h)
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()
	for i := 1; i <= r.replicas; i++ {
		h := newHash[i-1]
		// skip already exist node
		if r.contains(h) {
			return
		}

		r.sorted = append(r.sorted, h)
		r.vnodes[h] = node
	}
	sort.Slice(r.sorted, func(i, j int) bool { return r.sorted[i] < r.sorted[j] })
}

func (r *ring) contains(h uint32) bool {
	if _, exist := r.vnodes[h]; exist {
		return true
	}
	return false
}

func (r *ring) Remove(node Node) {
	ringVar.Remove.Add(1)

	newHash := make([]uint32, 0, r.replicas)
	for i := 1; i <= r.replicas; i++ {
		h := hash(node.String() + "_" + strconv.Itoa(i))
		newHash = append(newHash, h)
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()
	for i := 1; i <= r.replicas; i++ {
		h := newHash[i-1]
		// skip not exist node
		if !r.contains(h) {
			return
		}

		idx := sort.Search(len(r.sorted), func(i int) bool { return r.sorted[i] >= h })
		r.sorted = append(r.sorted[:idx], r.sorted[idx+1:]...)
		delete(r.vnodes, h)
	}
	sort.Slice(r.sorted, func(i, j int) bool { return r.sorted[i] < r.sorted[j] })
}

func (r *ring) Get(key string) Node {
	ringVar.Get.Add(1)

	if len(r.sorted) == 0 {
		return nil
	}
	h := hash(key)
	r.mtx.Lock()
	defer r.mtx.Unlock()
	idx := sort.Search(len(r.sorted), func(i int) bool { return r.sorted[i] >= h })

	// Means we have cycled back to the first replica.
	if idx == len(r.sorted) {
		idx = 0
	}

	return r.vnodes[r.sorted[idx]]
}

func hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

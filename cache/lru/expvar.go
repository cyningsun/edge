package lru

import (
	"expvar"
)

var m = struct {
	Get    *expvar.Int
	Set    *expvar.Int
	Delete *expvar.Int
	Exists *expvar.Int
	Hit    *expvar.Int
	Evict  *expvar.Int
}{
	Get:    expvar.NewInt("cache.lru.get"),
	Set:    expvar.NewInt("cache.lru.set"),
	Delete: expvar.NewInt("cache.lru.delete"),
	Exists: expvar.NewInt("cache.lru.exists"),
	Hit:    expvar.NewInt("cache.lru.hit"),
	Evict:  expvar.NewInt("cache.lru.evict"),
}

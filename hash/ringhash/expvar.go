package ringhash

import (
	"expvar"
)

var m = struct {
	Add    *expvar.Int
	Remove *expvar.Int
	Get    *expvar.Int
}{
	Add:    expvar.NewInt("hash.ringhash.add"),
	Remove: expvar.NewInt("hash.ringhash.remove"),
	Get:    expvar.NewInt("hash.ringhash.get"),
}

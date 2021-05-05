package edge

import (
	"time"
)

type Cache interface {
	Set(key string, val interface{}, d time.Duration)
	Get(key string) (value interface{}, ok bool)
	Delete(key string) (present bool)
	Exists(key string) bool
	Cap() int
	Len() int
}

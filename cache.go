package edge

// Cache is the common interface implemented by many kinds of cache algorithms
type Cache interface {
	// Set saves value under the key, treat the key as nonexistent
	// If key already exist, old value will be return
	// If key not exist, nil will be return
	Set(key string, val interface{}) interface{}

	// Get reads value under the key.
	// If key not exist, < nil, false> will be return
	Get(key string) (value interface{}, ok bool)

	// Delete removes the key in cache
	// If key not exist, false will be return
	Delete(key string) (present bool)

	// Exists returns whether the key exists in the cache.
	Exists(key string) bool

	// Cap returns maximum capacity of the cache.
	Cap() int

	// Len returns how many keys are currently stored in the cache.
	Len() int
}

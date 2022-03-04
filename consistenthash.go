package edge

// Server ConsistentHash uses String() as key when add into hash
type Server interface {
	String() string
}

// ConsistentHash was first described in a paper, ["Consistent hashing and random trees: Distributed caching protocols for relieving hot spots on the World Wide Web (1997)"](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.147.1879) by David Karger et al.
// It is used in distributed storage systems like Amazon Dynamo and memcached.
type ConsistentHash interface {
	// Add Saves server node into hash
	Add(s Server)

	// Remove deletes server node from hash
	Remove(node Server)

	// Get returns an server node close to key hash.
	Get(key string) Server
}

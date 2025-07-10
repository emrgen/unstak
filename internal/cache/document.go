package cache

// ObjectCache is a cache for generic object
// It is used to store published documents in memory to avoid fetching them from the database
// When a document is published, it is stored in the cache and when a document is unpublished, it is removed from the cache
// A document in a cache has a TTL to avoid storing it indefinitely
type ObjectCache[T any] struct {
	cache map[string]*T
}

// NewObjectCache creates a new ObjectCache
func NewObjectCache[T any]() *ObjectCache[T] {
	return &ObjectCache[T]{
		cache: make(map[string]*T),
	}
}

// GetDocument returns a document from the cache
func (c *ObjectCache[T]) GetDocument(key string) *T {
	return c.cache[key]
}

// SetDocument sets a document in the cache
func (c *ObjectCache[T]) SetDocument(key string, object *T) {
	c.cache[key] = object
}

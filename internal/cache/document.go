package cache

import v1 "github.com/emrgen/document/apis/v1"

// DocumentCache is a cache for documents
// It is used to store published documents in memory to avoid fetching them from the document service
// When a document is published, it is stored in the cache and when a document is unpublished, it is removed from the cache
// A document in a cache has a TTL to avoid storing it indefinitely
type DocumentCache struct {
	cache map[string]*v1.Document
}

// NewDocumentCache creates a new DocumentCache
func NewDocumentCache() *DocumentCache {
	return &DocumentCache{
		cache: make(map[string]*v1.Document),
	}
}

// GetDocument returns a document from the cache
func (c *DocumentCache) GetDocument(key string) *v1.Document {
	return c.cache[key]
}

// SetDocument sets a document in the cache
func (c *DocumentCache) SetDocument(key string, document *v1.Document) {
	c.cache[key] = document
}

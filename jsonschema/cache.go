package jsonschema

import (
	"fmt"
	"sync"

	"golang.org/x/sync/singleflight"
)

// Cache implements concurrency-safe caching with blocking load.
// (i.e., if the key is not cached then only one build/store happens.)
// The singleflight only blocks concurrent generation of the same cache
// key so syncronization is still needed on the storage map as multiple
// different keys can be accessed simultaneously.
type Cache[T any] struct {
	data   sync.Map
	flight singleflight.Group
}

func (c *Cache[T]) Load(key string, builder func() (T, error)) (T, error) {
	vCache, err, _ := c.flight.Do(key, func() (any, error) {
		if val, exists := c.data.Load(key); exists {
			// Cache hit.
			return val, nil
		}
		// Cache miss.
		v, err := builder()
		if err != nil {
			// Build fail.
			return v, err
		}
		// Build succes.
		c.data.Store(key, v)
		return v, nil
	})
	if err != nil {
		var zero T
		return zero, err
	}
	vReal, ok := vCache.(T)
	if !ok {
		return vReal, fmt.Errorf("invalid cached value; want %T, got %T", vReal, vCache)
	}
	return vReal, nil
}

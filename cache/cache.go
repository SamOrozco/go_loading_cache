package cache

import "time"

type CacheType int

const (
	// Blocking This cache type will block on an expired entry until the value is reloaded.
	Blocking CacheType = 0
	// Refresh This cache type will reload the value in the background when the entry is expired.
	Refresh CacheType = 1
)

type CacheInfo[K comparable, V any] struct {
	// MaxSize is the maximum number of entries that the cache can hold. If the cache already has more entries than the new maximum size, the cache will evict entries until the size is less than or equal to the new maximum size.
	MaxSize *int
	// Expiration is the expiration time for entries in the cache. If an entry is not accessed for longer than the expiration time, it will be evicted from the cache.
	Expiration time.Duration
	// CacheType will set the type of the cache to be implemented. For example a refresh cache will not block on an expired entry but reload in the background. See cache types.
	CacheType CacheType
	// CacheLoader is the loader that will be used to load values for the cache.
	CacheLoader CacheLoader[K, V]
	// EvictionPercent is the percent of the max size to delete when the cache is full
	EvictionPercent *int
	// Hooks are hooks that can be set on a cache to be called when certain events occur.
	Hooks CacheHooks[K]
}

// CacheHooks are hooks that can be set on a cache to be called when certain events occur.
type CacheHooks[K comparable] struct {
	OnCacheMiss         func(k K)
	OnCacheHit          func(k K)
	OnFailedToLoadEntry func(k K)
	OnCacheLoadDuration func(k K, duration time.Duration)
}

func (cacheInfo CacheInfo[K, V]) GetEvictionSize() int {
	return *cacheInfo.MaxSize * *cacheInfo.EvictionPercent / 100
}

// Load will be called by the cache when a key is not found in the cache. The loader should return the value associated with the key, or an error if the value could not be loaded.
type CacheLoader[K comparable, V any] func(k K) (V, error)

type CacheBuilder[K comparable, V any] interface {
	// SetMaxSize sets the maximum number of entries that the cache can hold. If the cache already has more entries than the new maximum size, the cache will evict entries until the size is less than or equal to the new maximum size.
	// Defaults to 100
	SetMaxSize(size int) CacheBuilder[K, V]
	// SetExpiration sets the expiration time for entries in the cache. If an entry is not accessed for longer than the expiration time, it will be evicted from the cache.
	// Defaults to 0 (no expiration)
	SetExpiration(expiration time.Duration) CacheBuilder[K, V]
	// SetCacheType will set the type of the cache to be implemented. For example a refresh cache will not block on an expired entry but reload in the background. See cache types.
	SetCacheType(cacheType CacheType) CacheBuilder[K, V]
	// SetEvictionPercent sets the percent of the max size to delete when the cache is full
	SetEvictionPercent(evictionPercent int) CacheBuilder[K, V]
	// Build creates a new cache with the specified loader and configuration.
	Build(loader CacheLoader[K, V]) Cache[K, V]
}

type Cache[K comparable, V any] interface {
	// Get returns the value associated with the key k from the cache. The second return type will be false if the key was not able to be loaded.
	// so a second return type of false will indicate an issue in loading the value and it will not be put in the cache
	Get(k K) (V, bool)
	// Put inserts a value into the cache associated with the key k. If the key already exists, the value will be updated.
	// The return value will be true if the value was inserted, and false if the value was updated.
	Put(k K, v V) bool
}

package cache

import "time"

const defaultMaxSize = 10
const defaultEvictionPercent = 10

type cacheBuilder[K comparable, V any] struct {
	cacheInfo    CacheInfo[K, V]
	cacheFactory CacheFactory[K, V]
}

func NewCacheBuilder[K comparable, V any]() CacheBuilder[K, V] {
	return &cacheBuilder[K, V]{
		cacheInfo:    CacheInfo[K, V]{},
		cacheFactory: NewCacheTypeFactory[K, V](),
	}
}

func NewCacheBuilderWithFactory[K comparable, V any](cacheFactory CacheFactory[K, V]) CacheBuilder[K, V] {
	return &cacheBuilder[K, V]{
		cacheInfo:    CacheInfo[K, V]{},
		cacheFactory: cacheFactory,
	}
}

func (c *cacheBuilder[K, V]) SetMaxSize(size int) CacheBuilder[K, V] {
	c.cacheInfo.MaxSize = &size
	return c
}

func (c *cacheBuilder[K, V]) SetExpiration(expiration time.Duration) CacheBuilder[K, V] {
	c.cacheInfo.Expiration = expiration
	return c
}

func (c *cacheBuilder[K, V]) SetEvictionPercent(evictionPercent int) CacheBuilder[K, V] {
	c.cacheInfo.EvictionPercent = &evictionPercent
	return c
}

func (c *cacheBuilder[K, V]) SetCacheType(cacheType CacheType) CacheBuilder[K, V] {
	c.cacheInfo.CacheType = cacheType
	return c
}

func (c *cacheBuilder[K, V]) Build(loader CacheLoader[K, V]) Cache[K, V] {

	if c.cacheInfo.MaxSize == nil {
		c.cacheInfo.MaxSize = PointerTo(defaultMaxSize)
	}

	if c.cacheInfo.EvictionPercent == nil {
		c.cacheInfo.EvictionPercent = PointerTo(defaultEvictionPercent)
	}

	c.cacheInfo.CacheLoader = loader
	return c.cacheFactory.BuildCache(c.cacheInfo)
}

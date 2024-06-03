package cache

type CacheFactory[K comparable, V any] interface {
	BuildCache(cacheInfo CacheInfo[K, V]) Cache[K, V]
}

type CacheTypeCacheFactory[K comparable, V any] struct {
	Clock Clock
}

func NewCacheTypeFactory[K comparable, V any]() CacheTypeCacheFactory[K, V] {
	return CacheTypeCacheFactory[K, V]{
		Clock: LocalClock{},
	}
}

func (c CacheTypeCacheFactory[K, V]) BuildCache(cacheInfo CacheInfo[K, V]) Cache[K, V] {
	switch cacheInfo.CacheType {
	case Refresh:
		return &refreshingExpiredCache[K, V]{
			cacheInfo: cacheInfo,
			cacheData: NewCacheData[K, V](c.Clock),
			clock:     c.Clock,
		}
	case Blocking:
		return &blockingExpiredCache[K, V]{
			cacheInfo: cacheInfo,
			cacheData: NewCacheData[K, V](c.Clock),
			clock:     c.Clock,
		}
	default:
		return &refreshingExpiredCache[K, V]{
			cacheInfo: cacheInfo,
			cacheData: NewCacheData[K, V](c.Clock),
			clock:     c.Clock,
		}
	}
}

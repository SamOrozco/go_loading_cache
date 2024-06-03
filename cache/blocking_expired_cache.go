package cache

import "time"

type blockingExpiredCache[K comparable, V any] struct {
	cacheInfo CacheInfo[K, V]
	cacheData CacheData[K, V]

	clock Clock
}

func (b *blockingExpiredCache[K, V]) Get(k K) (V, bool) {
	value, exists := b.cacheData.Get(k)
	if !exists || b.cacheData.IsExpired(k, b.cacheInfo.Expiration) {
		b.cacheMiss(k)

		value, err := b.loadCacheValue(k)
		if err != nil {
			b.failedToLoadEntry(k)
			return value, false
		}

		b.cacheData.Put(k, value)
		return value, true
	} else {
		b.cacheHit(k)
		return value, true
	}
}

func (b *blockingExpiredCache[K, V]) Put(k K, v V) bool {
	if b.cacheData.GetSize() >= *b.cacheInfo.MaxSize {
		b.cacheData.RemoveLeastRecentlyAccessed(b.cacheInfo.GetEvictionSize())
	}
	return b.cacheData.Put(k, v)
}

func (b *blockingExpiredCache[K, V]) cacheMiss(k K) {
	if b.cacheInfo.Hooks.OnCacheMiss != nil {
		b.cacheInfo.Hooks.OnCacheMiss(k)
	}
}

func (b *blockingExpiredCache[K, V]) cacheHit(k K) {
	if b.cacheInfo.Hooks.OnCacheHit != nil {
		b.cacheInfo.Hooks.OnCacheHit(k)
	}
}

func (b *blockingExpiredCache[K, V]) failedToLoadEntry(k K) {
	if b.cacheInfo.Hooks.OnFailedToLoadEntry != nil {
		b.cacheInfo.Hooks.OnFailedToLoadEntry(k)
	}
}

func (b *blockingExpiredCache[K, V]) loadCacheValue(k K) (V, error) {
	startLoad := b.clock.Now()
	value, err := b.cacheInfo.CacheLoader(k)
	if b.cacheInfo.Hooks.OnCacheLoadDuration != nil {
		b.cacheInfo.Hooks.OnCacheLoadDuration(k, time.Since(startLoad))
	}
	return value, err
}

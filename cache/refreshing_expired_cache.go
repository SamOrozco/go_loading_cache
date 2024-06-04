package cache

import "time"

type refreshingExpiredCache[K comparable, V any] struct {
	cacheInfo CacheInfo[K, V]
	cacheData CacheData[K, V]

	clock Clock
}

func (r refreshingExpiredCache[K, V]) Get(k K) (V, bool) {
	value, exists := r.cacheData.Get(k)
	// if the value does not exist we need to load it synchronously and put in cache
	// this should be the only time this cache will block
	if !exists {
		r.cacheMiss(k)
		value, err := r.cacheInfo.CacheLoader(k)
		if err != nil {
			r.failedToLoadEntry(k)
			return value, false
		}
		r.cacheData.Put(k, value)
		return value, true
	} else if r.cacheData.IsExpired(k, r.cacheInfo.Expiration) {
		r.cacheMiss(k)
		go func() {
			value, err := r.cacheInfo.CacheLoader(k)
			if err != nil {
				r.failedToLoadEntry(k)
				return
			}
			r.cacheData.Put(k, value)
		}()
	} else {
		r.cacheHit(k)
	}
	return value, true
}

func (r refreshingExpiredCache[K, V]) Put(k K, v V) bool {
	if r.cacheData.GetSize() >= *r.cacheInfo.MaxSize {
		r.cacheData.RemoveLeastRecentlyAccessed(r.cacheInfo.GetEvictionSize())
	}
	return r.cacheData.Put(k, v)
}

func (b refreshingExpiredCache[K, V]) Remove(k K) bool {
	b.cacheRemoved(k)
	return b.cacheData.Remove(k)
}

func (b refreshingExpiredCache[K, V]) cacheRemoved(k K) {
	if b.cacheInfo.Hooks.OnCacheRemove != nil {
		b.cacheInfo.Hooks.OnCacheRemove(k)
	}
}

func (r refreshingExpiredCache[K, V]) cacheMiss(k K) {
	if r.cacheInfo.Hooks.OnCacheMiss != nil {
		r.cacheInfo.Hooks.OnCacheMiss(k)
	}
}

func (r refreshingExpiredCache[K, V]) cacheHit(k K) {
	if r.cacheInfo.Hooks.OnCacheHit != nil {
		r.cacheInfo.Hooks.OnCacheHit(k)
	}
}

func (r refreshingExpiredCache[K, V]) failedToLoadEntry(k K) {
	if r.cacheInfo.Hooks.OnFailedToLoadEntry != nil {
		r.cacheInfo.Hooks.OnFailedToLoadEntry(k)
	}
}

func (r refreshingExpiredCache[K, V]) loadCacheValue(k K) (V, error) {
	startLoad := r.clock.Now()
	value, err := r.cacheInfo.CacheLoader(k)
	if r.cacheInfo.Hooks.OnCacheLoadDuration != nil {
		r.cacheInfo.Hooks.OnCacheLoadDuration(k, time.Since(startLoad))
	}
	return value, err
}

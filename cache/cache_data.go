package cache

import (
	"sync"
	"time"
)

type CacheData[K comparable, V any] interface {
	Cache[K, V]
	GetSize() int
	IsExpired(key K, cacheDuration time.Duration) bool
	RemoveLeastRecentlyAccessed(numToDelete int)
}

type cacheKey[K comparable] struct {
	// key value
	key K
	// insertTime the time the key was initially inserted into the cache
	insertTime time.Time
	// lastAccessTime the time the key was last accessed
	lastAccessTime time.Time
	// lastUpdateTime the time the key's value was last updated
	lastUpdateTime time.Time
}

// cacheData - thread safe data store
type cacheData[K comparable, V any] struct {
	keyData   map[K]*cacheKey[K]
	valueData map[K]V
	dataLock  *sync.Mutex
	clock     Clock
}

func NewCacheData[K comparable, V any](clockVar Clock) CacheData[K, V] {
	return &cacheData[K, V]{
		keyData:   make(map[K]*cacheKey[K]),
		valueData: make(map[K]V),
		dataLock:  &sync.Mutex{},
		clock:     clockVar,
	}

}

func (c *cacheData[K, V]) GetSize() int {
	return len(c.keyData)
}

func (c *cacheData[K, V]) IsExpired(key K, cacheDuration time.Duration) bool {
	// we are handling a special case if the duration is less than 1, that means an entry is never expired
	if cacheDuration < 1 {
		return false
	}

	cacheKey, exists := c.keyData[key]
	if !exists {
		return false
	} else {
		return c.clock.Now().After(cacheKey.lastUpdateTime.Add(cacheDuration))
	}
}

func (c *cacheData[K, V]) RemoveLeastRecentlyAccessed(numToDelete int) {

	// create a list that will maintain order by last accessed
	sortedList := LinkedSortedList[*cacheKey[K]]{
		CompareFunc: func(left *cacheKey[K], right *cacheKey[K]) int {
			return int(left.lastAccessTime.Sub(right.lastAccessTime).Nanoseconds())
		},
	}

	// add to sorted list
	for _, key := range c.keyData {
		sortedList.Add(key)
	}

	// remove the keys
	keysToDelete := sortedList.GetFirstN(numToDelete)
	for _, key := range keysToDelete {
		delete(c.keyData, key.key)
		delete(c.valueData, key.key)
	}
}

func (c *cacheData[K, V]) Get(k K) (V, bool) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()

	key, exists := c.keyData[k]
	if !exists {
		var defaultValue V
		return defaultValue, false
	} else {
		value := c.valueData[k]
		key.lastAccessTime = c.clock.Now()
		return value, true
	}
}

func (c *cacheData[K, V]) Put(k K, v V) bool {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()

	key, exists := c.keyData[k]
	if !exists {
		key = &cacheKey[K]{
			key:            k,
			insertTime:     c.clock.Now(),
			lastAccessTime: c.clock.Now(),
			lastUpdateTime: c.clock.Now(),
		}
		c.keyData[k] = key
		c.valueData[k] = v
		return true
	} else {
		key.lastUpdateTime = c.clock.Now()
		c.valueData[k] = v
		return false
	}
}

func (c *cacheData[K, V]) Remove(k K) bool {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()

	_, exists := c.keyData[k]
	delete(c.keyData, k)
	delete(c.valueData, k)
	return exists
}

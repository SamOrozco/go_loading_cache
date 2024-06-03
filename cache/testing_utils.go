package cache

import (
	"time"
)

type TestCacheLoader[K comparable, V any] struct {
	ReturnValues   []V
	KeysRequests   []K
	returnValueIdx int
}

func (t *TestCacheLoader[K, V]) Load(k K) (V, error) {
	t.KeysRequests = append(t.KeysRequests, k)
	if len(t.ReturnValues) == 0 {
		var defaultValue V
		return defaultValue, nil
	}

	if t.returnValueIdx >= len(t.ReturnValues) {
		t.returnValueIdx = 0
	}
	returnValue := t.ReturnValues[t.returnValueIdx]
	t.returnValueIdx++
	return returnValue, nil
}

func BuildTestCacheByType[K comparable, V any](cacheType CacheType, loader CacheLoader[K, V], clock Clock) Cache[K, V] {
	return NewCacheBuilderWithFactory[K, V](CacheTypeCacheFactory[K, V]{
		Clock: clock,
	}).
		SetCacheType(cacheType).
		Build(loader)
}

func BuildTestCacheByTypeAndExpirationMillis[K comparable, V any](cacheType CacheType, loader CacheLoader[K, V], clock Clock, expirationMillis int) Cache[K, V] {
	return NewCacheBuilderWithFactory[K, V](CacheTypeCacheFactory[K, V]{
		Clock: clock,
	}).
		SetCacheType(cacheType).
		SetExpiration(time.Millisecond * time.Duration(expirationMillis)).
		Build(loader)
}

type TestClock struct {
	times []time.Time
	idx   int
}

func NewTestClock(times ...time.Time) Clock {
	return &TestClock{
		times: times,
	}
}

func (t *TestClock) Now() time.Time {
	if t.idx >= len(t.times) {
		t.idx = 0
	}
	time := t.times[t.idx]
	t.idx++
	return time
}

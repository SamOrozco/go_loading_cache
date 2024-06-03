package cache

import "testing"

type TestCacheFactory[K comparable, V any] struct {
	cacheInfo CacheInfo[K, V]
}

func (t *TestCacheFactory[K, V]) BuildCache(cacheInfo CacheInfo[K, V]) Cache[K, V] {
	t.cacheInfo = cacheInfo
	return nil
}

func TestUsingBuilderWillSetTheProperValuesOnCache(t *testing.T) {
	// setup
	testCacheFactory := &TestCacheFactory[string, string]{}
	cb := NewCacheBuilderWithFactory[string, string](testCacheFactory)

	// execute
	cb.SetMaxSize(10).
		SetExpiration(10).
		SetCacheType(Blocking).
		SetEvictionPercent(10).
		Build(func(k string) (string, error) {
			return "value", nil
		})
	// verify
	if *testCacheFactory.cacheInfo.MaxSize != 10 {
		t.Errorf("MaxSize not set properly")
	}
	if testCacheFactory.cacheInfo.Expiration != 10 {
		t.Errorf("Expiration not set properly")
	}
	if testCacheFactory.cacheInfo.CacheType != Blocking {
		t.Errorf("CacheType not set properly")
	}
	if *testCacheFactory.cacheInfo.EvictionPercent != 10 {
		t.Errorf("EvictionPercent not set properly")
	}
}

func TestBuilderWillUseCorrectDefaultValuesIfNoneSet(t *testing.T) {
	// setup
	testCacheFactory := &TestCacheFactory[string, string]{}
	cb := NewCacheBuilderWithFactory[string, string](testCacheFactory)

	// execute
	cb.Build(func(k string) (string, error) {
		return "value", nil
	})
	// verify
	if *testCacheFactory.cacheInfo.MaxSize != 10 {
		t.Errorf("MaxSize not set properly")
	}
	if testCacheFactory.cacheInfo.Expiration != 0 {
		t.Errorf("Expiration not set properly")
	}
	if testCacheFactory.cacheInfo.CacheType != Blocking {
		t.Errorf("CacheType not set properly")
	}
	if *testCacheFactory.cacheInfo.EvictionPercent != 10 {
		t.Errorf("EvictionPercent not set properly")
	}
}

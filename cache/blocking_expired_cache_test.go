package cache

import (
	"testing"
	"time"
)

var cache Cache[string, string]
var cacheLoader TestCacheLoader[string, string]

func initTests() {
	cacheLoader = TestCacheLoader[string, string]{
		ReturnValues: []string{"value", "second"},
	}
	cache = BuildTestCacheByType[string, string](Blocking, cacheLoader.Load, LocalClock{})
}

func TestWhenCacheIsEmptyWillLoadDataAndPutInCache(t *testing.T) {
	// setup
	initTests()

	// executor
	keyValue, wasSet := cache.Get("key")

	// verify
	if !wasSet {
		t.Errorf("Expected value to exist")
	}
	if keyValue != "value" {
		t.Errorf("Expected value to be 'value'")
	}
	if len(cacheLoader.KeysRequests) != 1 {
		t.Errorf("Expected to call loader once")
	}
	if cacheLoader.KeysRequests[0] != "key" {
		t.Errorf("Expected to call loader with key 'key'")
	}
}

func TestWhenLoadingAValueThatIsInTheCacheAndNotExpiredDoNotLoadData(t *testing.T) {
	// setup
	initTests()

	// execute
	firstValue, wasSet := cache.Get("key")
	secondValue, wasSet2 := cache.Get("key")

	// verify
	if !wasSet {
		t.Errorf("Expected value to exist")
	}

	if !wasSet2 {
		t.Errorf("Expected value to exist")
	}

	if firstValue != "value" {
		t.Errorf("Expected value to be 'value'")
	}

	if secondValue != "value" {
		t.Errorf("Expected value to be 'value'")
	}

	if len(cacheLoader.KeysRequests) != 1 {
		t.Errorf("Expected to call loader once")
	}

	if cacheLoader.KeysRequests[0] != "key" {
		t.Errorf("Expected to call loader with key 'key'")
	}
}

func TestWhenLoadingDataThatHasExpiredReloadDataOnRequest(t *testing.T) {
	// setup
	initTests()
	testCacheTimes := []time.Time{
		time.Unix(1000, 0), // load cache value - start load - in case we are timing the cache request - not important for test
		time.Unix(1000, 0), // initial insert - insertTime - not important for test
		time.Unix(1000, 0), // initial insert - lastAccessTime - not important for test
		time.Unix(1000, 0), // initial insert - lastUpdateTime - important for test
		time.Unix(1000, 0), // second get request - value exists update last accessed time - not important for test
		time.Unix(2000, 0), // second get request - checking to see if the entry is expired - important for test - should expire initial request
	}
	cache = BuildTestCacheByTypeAndExpirationMillis[string, string](Blocking, cacheLoader.Load, NewTestClock(testCacheTimes...), 10)

	// execute
	value1, wasLoaded := cache.Get("key")
	value2, wasLoaded2 := cache.Get("key")

	// verify
	if !wasLoaded {
		t.Errorf("Expected value to be loaded")
	}

	if !wasLoaded2 {
		t.Errorf("Expected value to be loaded")
	}

	if value1 != "value" {
		t.Errorf("Expected value to be 'value'")
	}

	if value2 != "second" {
		t.Errorf("Expected value to be 'second'")
	}

	if len(cacheLoader.KeysRequests) != 2 {
		t.Errorf("Expected to call loader twice")
	}
}

func TestWhenRemovingFromCacheAndItemExistsReturnTrue(t *testing.T) {
	// setup
	initTests()

	// execute
	cache.Put("key", "value")
	removed := cache.Remove("key")

	// verify
	if !removed {
		t.Errorf("Expected to remove value")
	}
}

func TestWhenRemovingFromCacheAndItemDoesNotExistReturnFalse(t *testing.T) {
	// setup
	initTests()

	// execute
	removed := cache.Remove("key")

	// verify
	if removed {
		t.Errorf("Expected value to not exists when removing")
	}
}

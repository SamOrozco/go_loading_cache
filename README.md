# Go Loading cache

This is a simple straightforward cache library for Go. It is takes inspiration from the [Caffeine](https://github.com/ben-manes/caffeine) cache library for Java.

## Installing
```bash
go get -u github.com/SamOrozco/go_loading_cache
```

## Features
There are two main caches types in this library:
1. BlockingCache: This cache is a blocking cache that will block on a get request in the entry is expired until the value is available.
2. RefreshingCache: This cache is a cache that will refresh the value in a background when it is expired. This allows for requests to cache to be non-blocking.

## Usage


### Cache interface
```go
type Cache[K comparable, V any] interface {
	// Get returns the value associated with the key k from the cache. The second return type will be false if the key was not able to be loaded.
	// so a second return type of false will indicate an issue in loading the value and it will not be put in the cache
	Get(k K) (V, bool)
	// Put inserts a value into the cache associated with the key k. If the key already exists, the value will be updated.
	// The return value will be true if the value was inserted, and false if the value was updated.
	Put(k K, v V) bool
}
```

### Using the cache builder
```go
	userCache := cache.NewCacheBuilder[string, *User]().
		SetMaxSize(100). // max number of items in the cache before we remove items - default 10
		SetExpiration(time.Second * 0). // entry expiration duration - default 0 or no expiration
		SetCacheType(cache.Blocking). // cache type - refreshing or blocking - default blocking
		SetEvictionPercent(10). // percent of the max size to delete when the cache is full - default 10
		Build(func(k string) (*User, error) {
			// loading function
			// this is used to get the value for the given key and insert into the cache
			return &User{}, nil
		})
```


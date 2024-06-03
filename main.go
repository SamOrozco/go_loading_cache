package main

import (
	"go_loading_cache/cache"
	"time"
)

type User struct {
}

func main() {
	userCache := cache.NewCacheBuilder[string, *User]().
		SetMaxSize(100).                // max number of items in the cache before we remove items - default 10
		SetExpiration(time.Second * 0). // entry expiration duration - default 0 or no expiration
		SetCacheType(cache.Blocking).   // cache type - refreshing or blocking - default blocking
		SetEvictionPercent(10).         // percent of the max size to delete when the cache is full - default 10
		Build(func(k string) (*User, error) {
			// loading function
			// this is used to get the value for the given key and insert into the cache
			return &User{}, nil
		})

	userCache.Put("key", &User{})
}

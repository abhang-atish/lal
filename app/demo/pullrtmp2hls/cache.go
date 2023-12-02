package main

import (
	"github.com/coocood/freecache"
)

var cache *freecache.Cache

func InitCache() {
	cacheSize := 100 * 1024 * 1024
	cache = freecache.NewCache(cacheSize)
}

func GetValue(key []byte) ([]byte, error) {
	return cache.Get(key)
}

func SetValue(key []byte, value []byte, expiry int) {
	err := cache.Set(key, value, expiry)
	if err != nil {

	}
}

func DeleteKey(key []byte) bool {
	return cache.Del(key)
}

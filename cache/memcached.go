package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

var (
	memcacheClient *memcache.Client
)

func init() {
	memcacheClient = memcache.New("localhost:11211")
}

func CacheMemcacheSet(key string, value string) error {
	return memcacheClient.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: 0,
	})
}

func CacheMemcacheGet(key string) (string, error) {
	item, err := memcacheClient.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

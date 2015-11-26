package cache

import ()

// Cache interface
// usage:
//	c.Put("key",value)
//	v := c.Get("key")
//
type LoadingCache interface {
	// get cached value by key.
	Get(key string) interface{}
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) []interface{}
	// set cached value with key
	Modify(key string, val interface{}) error
	// invalid cached value by key.
	Invalid(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// invalid all cache.
	InvalidAll() error
	//start
	StartAndGC() error
	//stop
	Stop()
}

type Loader interface {
	Load(key string) interface{}
	Modify(key string, value interface{}) error
}

func NewCache(loader Loader, refreshInterval int64) (lc LoadingCache, err error) {
	lc = NewMemoryCache(loader, refreshInterval)
	err = lc.StartAndGC()
	return
}

package cache

import (
	"github.com/HongdianLab/concurrent-map"

	"fmt"
	"sync"
	"time"
)

var (
	DefaultExpireAfterWrite  int64 = 86400 * 30 // 1 month
	DefaultExpireAfterAccess int64 = 86400 * 30 // 1 month
)

// Memory cache item.
type MemoryItem struct {
	val        interface{}
	Lastaccess time.Time
	Lastwrite  time.Time
}

// Memory cache contains a RW locker for safe map storage.
type MemoryCache struct {
	lock              sync.RWMutex
	dur               time.Duration
	items             cmap.ConcurrentMap // ConcurrentMap:map[string]*MemoryItem
	Every             int64              // run an refreshKey Every clock time
	loader            Loader
	expireAfterWrite  int64
	expireAfterAccess int64
	stop              chan bool
}

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache(loader Loader, refreshInterval, writeTimeout int64) *MemoryCache {
	if writeTimeout < 0 {
		writeTimeout = DefaultExpireAfterWrite
	}

	cmap.SHARD_COUNT = 128
	cache := MemoryCache{
		items:             cmap.New(),
		Every:             refreshInterval,
		dur:               time.Duration(refreshInterval) * time.Second,
		loader:            loader,
		expireAfterWrite:  writeTimeout,
		expireAfterAccess: DefaultExpireAfterAccess,
		stop:              make(chan bool, 1),
	}
	return &cache
}

func (bc *MemoryCache) isExpireAfterAccess() bool {
	return bc.expireAfterAccess > 0
}
func (bc *MemoryCache) isExpireAfterWrite() bool {
	return bc.expireAfterWrite > 0
}

func (bc *MemoryCache) isExpired(item *MemoryItem, now int64) bool {
	if bc.isExpireAfterAccess() && (now-item.Lastaccess.Unix()) > bc.expireAfterAccess {
		return true
	}
	if bc.isExpireAfterWrite() && (now-item.Lastwrite.Unix()) > bc.expireAfterWrite {
		return true
	}
	return false
}

// Get cache from memory.
// if non-existed or expired, return nil.
func (bc *MemoryCache) Get(name string) interface{} {
	if v, ok := bc.items.Get(name); ok {
		now := time.Now()
		item := v.(*MemoryItem)
		if bc.isExpired(item, now.Unix()) {
			go bc.Invalid(name)
			return nil
		}
		item.Lastaccess = now
		return item.val
	} else {
		go bc.put(name, nil)
	}
	return nil
}

// GetMulti gets caches from memory.
// if non-existed or expired, return nil.
func (bc *MemoryCache) GetMulti(names []string) []interface{} {
	var rc []interface{}
	for _, name := range names {
		rc = append(rc, bc.Get(name))
	}
	return rc
}

// modify value with name.
func (bc *MemoryCache) Modify(name string, value interface{}) error {
	err := bc.loader.Modify(name, value)
	if err != nil {
		return err
	}
	val := bc.loader.Load(name)

	bc.put(name, val)
	return err
}

func (bc *MemoryCache) put(name string, value interface{}) {
	now := time.Now()
	bc.items.Set(name, &MemoryItem{
		val:        value,
		Lastaccess: now,
		Lastwrite:  now,
	})
}

/// Invalid cache in memory.
func (bc *MemoryCache) Invalid(name string) error {
	bc.items.Remove(name)
	return nil
}

// check cache exist in memory.
func (bc *MemoryCache) IsExist(name string) bool {
	ok := bc.items.Has(name)
	return ok
}

// delete all cache in memory.
func (bc *MemoryCache) InvalidAll() error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	bc.items = cmap.New()
	return nil
}

// start memory cache. it will check expiration in every clock time.
func (bc *MemoryCache) StartAndGC() error {
	go bc.vaccuum()
	return nil
}

func (bc *MemoryCache) Stop() {
	bc.stop <- true
}

// refresh and check expiration.
func (bc *MemoryCache) vaccuum() {
	if bc.Every < 1 {
		return
	}
	fmt.Printf("start refresh %v\n", bc.dur)
	heartbeatTicker := time.NewTicker(bc.dur)

	for {
		select {
		case <-heartbeatTicker.C:
			if bc.items == nil {
				return
			}
			tasks := bc.items.IterBuffered()

			bc.lock.Lock()
			var wg sync.WaitGroup
			for i := 0; i < 32; i++ {
				wg.Add(1)
				go func(bc *MemoryCache) {
					defer wg.Done()
					for tuple := range tasks {
						bc.refreshByName(tuple.Key)
					}
				}(bc)
			}
			wg.Wait()
			bc.lock.Unlock()
		case <-bc.stop:
			return
		}
	}
}

// refreshByName returns true if an item is expired.
func (bc *MemoryCache) refreshByName(name string) bool {
	v, ok := bc.items.Get(name)
	if !ok {
		return true
	}
	item := v.(*MemoryItem)
	if bc.isExpired(item, time.Now().Unix()) {
		bc.items.Remove(name)
		return true
	}
	val := bc.loader.Load(name)
	if val != nil {
		bc.put(name, val)
		return false
	} else {
		bc.items.Remove(name)
		return true
	}
}

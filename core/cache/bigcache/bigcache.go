package bigcache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"log"
)

var cache *bigcache.BigCache

func Init(opts ...Option) {
	options := NewOptions(opts...)
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: options.Shards,

		// time after which entry can be evicted
		LifeWindow: options.LifeWindow,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: options.CleanWindow,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: options.MaxEntriesInWindow,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: options.MaxEntrySize,

		// prints information about additional memory allocation
		Verbose: options.Verbose,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: options.HardMaxCacheSize,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: options.OnRemove,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
	}
	bc, err := bigcache.New(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	cache = bc
}

func BigCacheSet(key string, value string) {
	err := cache.Set(key, []byte(value))
	if err != nil {
		log.Panicln("bigcache set error")
	}
}

func BigCacheGet(key string) string {
	val, _ := cache.Get(key)
	return string(val)
}

func BigCacheDelete(key string) error {
	return cache.Delete(key)
}

func BigCacheClose() {
	cache.Close()
}

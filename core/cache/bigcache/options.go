package bigcache

import "time"

// Options BigCache配置结构体
type Options struct {
	// Number of cache shards, value must be a power of two
	Shards int
	// Time after which entry can be evicted
	LifeWindow time.Duration
	// Interval between removing expired entries (clean up).
	// If set to <= 0 then no action is performed. Setting to < 1 second is counterproductive — bigcache has a one second resolution.
	CleanWindow time.Duration
	// Max number of entries in life window. Used only to calculate initial size for cache shards.
	// When proper value is set then additional memory allocation does not occur.
	MaxEntriesInWindow int
	// Max size of entry in bytes. Used only to calculate initial size for cache shards.
	MaxEntrySize int
	// StatsEnabled if true calculate the number of times a cached resource was requested.
	StatsEnabled bool
	// Verbose mode prints information about new memory allocation
	Verbose bool
	// HardMaxCacheSize is a limit for BytesQueue size in MB.
	// It can protect application from consuming all available memory on machine, therefore from running OOM Killer.
	// Default value is 0 which means unlimited size. When the limit is higher than 0 and reached then
	// the oldest entries are overridden for the new ones. The max memory consumption will be bigger than
	// HardMaxCacheSize due to Shards' s additional memory. Every Shard consumes additional memory for map of keys
	// and statistics (map[uint64]uint32) the size of this map is equal to number of entries in
	// cache ~ 2×(64+32)×n bits + overhead or map itself.
	HardMaxCacheSize int
	// OnRemove is a callback fired when the oldest entry is removed because of its expiration time or no space left
	// for the new entry, or because delete was called.
	// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
	// ignored if OnRemoveWithMetadata is specified.
	OnRemove func(key string, entry []byte)
}

type Option func(options *Options)

func WithShards(shards int) Option {
	return func(o *Options) {
		o.Shards = shards
	}
}

func WithLifeWindow(lifeWindow time.Duration) Option {
	return func(o *Options) {
		o.LifeWindow = lifeWindow
	}
}

func WithCleanWindow(cleanWindow time.Duration) Option {
	return func(o *Options) {
		o.CleanWindow = cleanWindow
	}
}

func WithMaxEntriesInWindow(maxEntriesInWindow int) Option {
	return func(o *Options) {
		o.MaxEntriesInWindow = maxEntriesInWindow
	}
}

func WithMaxEntrySize(maxEntrySize int) Option {
	return func(o *Options) {
		o.MaxEntrySize = maxEntrySize
	}
}

func WithStatsEnabled(statsEnabled bool) Option {
	return func(o *Options) {
		o.StatsEnabled = statsEnabled
	}
}

func WithVerbose(verbose bool) Option {
	return func(o *Options) {
		o.Verbose = verbose
	}
}

func WithHardMaxCacheSize(hardMaxCacheSize int) Option {
	return func(o *Options) {
		o.HardMaxCacheSize = hardMaxCacheSize
	}
}

func WithOnRemove(onRemove func(key string, entry []byte)) Option {
	return func(o *Options) {
		o.OnRemove = onRemove
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Shards:             1024,
		LifeWindow:         10 * time.Minute,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           nil,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

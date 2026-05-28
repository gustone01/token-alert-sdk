package core

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	Enabled           bool
	Platform          Platform
	ServiceName       string
	HostIP            string
	PeekBytes         int
	DedupTTL          time.Duration
	NotifyTimeout     time.Duration
	EnrichFunc        EnrichFunc
	Sender            Sender
	Redis             *redis.Client
	PlatformEnabledFn func(Platform) bool
}

type Option func(*Options)

func WithEnabled(enabled bool) Option {
	return func(o *Options) { o.Enabled = enabled }
}

func WithPlatform(platform Platform) Option {
	return func(o *Options) { o.Platform = platform }
}

func WithServiceName(name string) Option {
	return func(o *Options) { o.ServiceName = name }
}

func WithHostIP(ip string) Option {
	return func(o *Options) { o.HostIP = ip }
}

func WithPeekBytes(n int) Option {
	return func(o *Options) { o.PeekBytes = n }
}

func WithDedupTTL(ttl time.Duration) Option {
	return func(o *Options) { o.DedupTTL = ttl }
}

func WithDedupRedis(rdb *redis.Client, ttl time.Duration) Option {
	return func(o *Options) {
		o.Redis = rdb
		if ttl > 0 {
			o.DedupTTL = ttl
		}
	}
}

func WithNotifyTimeout(d time.Duration) Option {
	return func(o *Options) { o.NotifyTimeout = d }
}

func WithEnrichFunc(fn EnrichFunc) Option {
	return func(o *Options) { o.EnrichFunc = fn }
}

func WithSender(sender Sender) Option {
	return func(o *Options) { o.Sender = sender }
}

func WithConfig(cfg Config) Option {
	return func(o *Options) {
		cfg.ApplyDefaults()
		o.Enabled = cfg.Enabled
		o.HostIP = cfg.HostIP
		o.PeekBytes = cfg.PeekBytes
		o.DedupTTL = cfg.DedupTTL
		if cfg.NotifyTimeoutSecs > 0 {
			o.NotifyTimeout = time.Duration(cfg.NotifyTimeoutSecs) * time.Second
		}
		p := cfg.Platforms
		o.PlatformEnabledFn = func(platform Platform) bool {
			if !cfg.Enabled {
				return false
			}
			switch platform {
			case PlatformBytedance:
				return p.Bytedance
			case PlatformTencent:
				return p.Tencent
			default:
				return false
			}
		}
	}
}

func applyOptions(opts []Option) Options {
	o := Options{
		PeekBytes:     DefaultPeekBytes,
		DedupTTL:      DefaultDedupTTL,
		NotifyTimeout: DefaultNotifyTimeoutSecs * time.Second,
	}
	for _, fn := range opts {
		fn(&o)
	}
	if o.PeekBytes <= 0 {
		o.PeekBytes = DefaultPeekBytes
	}
	if o.DedupTTL <= 0 {
		o.DedupTTL = DefaultDedupTTL
	}
	if o.NotifyTimeout <= 0 {
		o.NotifyTimeout = DefaultNotifyTimeoutSecs * time.Second
	}
	return o
}

func OptionsFromConfig(cfg Config, rdb *redis.Client, serviceName string, sender Sender, enrich EnrichFunc) []Option {
	cfg.ApplyDefaults()
	return []Option{
		WithConfig(cfg),
		WithServiceName(serviceName),
		WithDedupRedis(rdb, cfg.DedupTTL),
		WithSender(sender),
		WithEnrichFunc(enrich),
	}
}

var _ http.RoundTripper = (*Transport)(nil)

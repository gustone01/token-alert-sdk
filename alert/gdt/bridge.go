package gdt

import (
	"net/http"

	alertcore "github.com/gustone01/token-alert-sdk/alert/core"
	gdtlimiter "github.com/gustone01/oe-limiter-sdk/limiter/gdt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewTransport(db *gorm.DB, rdb *redis.Client, alertOpts []alertcore.Option, gdtOpts ...gdtlimiter.Option) (*gdtlimiter.Transport, error) {
	opts := append([]alertcore.Option{alertcore.WithPlatform(alertcore.PlatformTencent)}, alertOpts...)
	handler := alertcore.NewHandler(opts...)
	alertcore.RegisterHandler(alertcore.PlatformTencent, handler)

	alertTP := alertcore.NewTransport(http.DefaultTransport, opts...)
	allOpts := append(gdtOpts, gdtlimiter.WithBaseTransport(alertTP))
	return gdtlimiter.NewTransport(db, rdb, allOpts...)
}

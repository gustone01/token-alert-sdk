package oe

import (
	"net/http"

	alertcore "github.com/gustone01/token-alert-sdk/alert/core"
	oelimiter "github.com/gustone01/oe-limiter-sdk/limiter/oe"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewTransport(db *gorm.DB, rdb *redis.Client, alertOpts []alertcore.Option, oeOpts ...oelimiter.Option) (*oelimiter.Transport, error) {
	opts := append([]alertcore.Option{alertcore.WithPlatform(alertcore.PlatformBytedance)}, alertOpts...)
	handler := alertcore.NewHandler(opts...)
	alertcore.RegisterHandler(alertcore.PlatformBytedance, handler)

	alertTP := alertcore.NewTransport(http.DefaultTransport, opts...)
	allOpts := append(oeOpts, oelimiter.WithBaseTransport(alertTP))
	return oelimiter.NewTransport(db, rdb, allOpts...)
}

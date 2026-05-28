# token-alert-sdk

媒体平台（巨量引擎 / 腾讯广告）`access_token` 失效检测与企微告警插件，作为 [oe-limiter-sdk](https://github.com/gustone01/oe-limiter-sdk) 的 **Base Transport** 挂载。

## 安装

```bash
go get github.com/gustone01/token-alert-sdk@v0.1.0
```

## 快速接入（巨量 + 限流）

```go
import (
    alertcore "github.com/gustone01/token-alert-sdk/alert/core"
    alertoe "github.com/gustone01/token-alert-sdk/alert/oe"
    "github.com/gustone01/token-alert-sdk/alert/wecom"
    oelimiter "github.com/gustone01/oe-limiter-sdk/limiter/oe"
)

cfg := alertcore.Config{
    Enabled:         true,
    WecomWebhookKey: "your-webhook-key",
    DedupTTL:        30 * time.Minute,
    Platforms:       alertcore.PlatformSwitch{Bytedance: true, Tencent: true},
}
cfg.ApplyDefaults()

tp, err := alertoe.NewTransport(db, rdb,
    alertcore.WithConfig(cfg),
    alertcore.WithServiceName("my-service"),
    alertcore.WithSender(wecom.New(cfg.WecomWebhookKey, nil)),
    alertcore.WithDedupRedis(rdb, cfg.DedupTTL),
    oelimiter.WithOnDiscover(onDiscover),
)
```

## 腾讯 GDT

使用 `alert/gdt.NewTransport`，用法与 `alert/oe` 相同。

## 手动兜底

未走 Transport 时可调用：

```go
alertcore.TryNotify(alertcore.PlatformTencent, code, message, apiPath, accessToken)
```

## 配置说明

| 字段 | 说明 |
|------|------|
| `enabled` | 总开关，`false` 时零开销透传 |
| `wecom_webhook_key` | 企微机器人 key 或完整 webhook URL |
| `host_ip` | 告警来源机器 IP；为空时自动探测本机 IPv4 |
| `dedup_ttl` | 同一 token+错误码告警冷却（Redis SET NX） |
| `peek_bytes` | 响应 peek 上限，默认 8192 |

## Token 错误码

- 巨量：`40102` `40103` `40104`
- 腾讯：`11000` `11002` `11004` `11012`

## 本地联调

```text
gy_server/go.work:
use ../token-alert-sdk
```

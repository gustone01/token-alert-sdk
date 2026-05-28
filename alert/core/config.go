package core

import "time"

const (
	DefaultPeekBytes         = 8192
	DefaultDedupTTL          = 30 * time.Minute
	DefaultNotifyTimeoutSecs = 5
)

type PlatformSwitch struct {
	Bytedance bool `yaml:"bytedance" json:"bytedance"`
	Tencent   bool `yaml:"tencent" json:"tencent"`
}

type Config struct {
	Enabled           bool           `yaml:"enabled" json:"enabled"`
	WecomWebhookKey   string         `yaml:"wecom_webhook_key" json:"wecom_webhook_key"`
	HostIP            string         `yaml:"host_ip" json:"host_ip"`
	MentionedList     []string       `yaml:"mentioned_list" json:"mentioned_list"`
	DedupTTL          time.Duration  `yaml:"dedup_ttl" json:"dedup_ttl"`
	PeekBytes         int            `yaml:"peek_bytes" json:"peek_bytes"`
	Platforms         PlatformSwitch `yaml:"platforms" json:"platforms"`
	NotifyTimeoutSecs int            `yaml:"notify_timeout_secs" json:"notify_timeout_secs"`
}

func (c *Config) ApplyDefaults() {
	if c == nil {
		return
	}
	if c.PeekBytes <= 0 {
		c.PeekBytes = DefaultPeekBytes
	}
	if c.DedupTTL <= 0 {
		c.DedupTTL = DefaultDedupTTL
	}
	if c.NotifyTimeoutSecs <= 0 {
		c.NotifyTimeoutSecs = DefaultNotifyTimeoutSecs
	}
	if !c.Platforms.Bytedance && !c.Platforms.Tencent {
		c.Platforms.Bytedance = true
		c.Platforms.Tencent = true
	}
}

func (c *Config) PlatformEnabled(p Platform) bool {
	if c == nil || !c.Enabled {
		return false
	}
	switch p {
	case PlatformBytedance:
		return c.Platforms.Bytedance
	case PlatformTencent:
		return c.Platforms.Tencent
	default:
		return false
	}
}

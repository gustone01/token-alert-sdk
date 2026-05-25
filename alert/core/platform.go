package core

type Platform string

const (
	PlatformBytedance Platform = "bytedance"
	PlatformTencent   Platform = "tencent"
)

func (p Platform) DisplayName() string {
	switch p {
	case PlatformBytedance:
		return "巨量引擎"
	case PlatformTencent:
		return "腾讯广告"
	default:
		return string(p)
	}
}

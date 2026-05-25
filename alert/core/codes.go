package core

var (
	bytedanceTokenCodes = map[int64]struct{}{
		40102: {},
		40103: {},
		40104: {},
		40105: {},
	}
	tencentTokenCodes = map[int64]struct{}{
		11000: {},
		11002: {},
		11004: {},
		11012: {},
	}
)

func IsTokenError(platform Platform, code int64) bool {
	if code == 0 {
		return false
	}
	switch platform {
	case PlatformBytedance:
		_, ok := bytedanceTokenCodes[code]
		return ok
	case PlatformTencent:
		_, ok := tencentTokenCodes[code]
		return ok
	default:
		return false
	}
}

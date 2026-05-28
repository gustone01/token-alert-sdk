package wecom

import (
	"strings"
	"testing"
	"time"

	"github.com/gustone01/token-alert-sdk/alert/core"
)

func TestRenderContentIncludesIP(t *testing.T) {
	content := renderContent(core.Event{
		Service:    "my-service",
		HostIP:     "192.168.1.10",
		Platform:   core.PlatformBytedance,
		APIPath:    "/open_api/v1.0/test",
		Code:       40102,
		Message:    "token expired",
		OccurredAt: time.Date(2026, 5, 28, 12, 0, 0, 0, time.UTC),
	})
	if !strings.Contains(content, "IP：192.168.1.10") {
		t.Fatalf("content missing ip: %s", content)
	}
}

func TestRenderContentUnknownIP(t *testing.T) {
	content := renderContent(core.Event{
		Service:  "my-service",
		Platform: core.PlatformTencent,
		Code:     11000,
	})
	if !strings.Contains(content, "IP：未知") {
		t.Fatalf("content missing unknown ip: %s", content)
	}
}

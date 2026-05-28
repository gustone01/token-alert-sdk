package core_test

import (
	"testing"

	"github.com/gustone01/token-alert-sdk/alert/core"
)

func TestResolveHostIPConfigured(t *testing.T) {
	if got := core.ResolveHostIP("10.1.2.3"); got != "10.1.2.3" {
		t.Fatalf("unexpected ip: %s", got)
	}
}

func TestResolveHostIPAutoDetect(t *testing.T) {
	if got := core.ResolveHostIP(""); got == "" {
		t.Fatal("expected auto-detected ipv4")
	}
}

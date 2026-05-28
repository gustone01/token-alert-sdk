package core_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gustone01/token-alert-sdk/alert/core"
)

func TestIsTokenError(t *testing.T) {
	if !core.IsTokenError(core.PlatformBytedance, 40102) {
		t.Fatal("expected 40102 token error")
	}
	if core.IsTokenError(core.PlatformBytedance, 0) {
		t.Fatal("code 0 should not be token error")
	}
	if !core.IsTokenError(core.PlatformTencent, 11000) {
		t.Fatal("expected 11000 token error")
	}
}

func TestTransportPeekPreservesBody(t *testing.T) {
	body := `{"code":40102,"message":"token expired","data":{}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, body)
	}))
	defer srv.Close()

	var notified bool
	done := make(chan struct{})
	tp := core.NewTransport(http.DefaultTransport,
		core.WithEnabled(true),
		core.WithPlatform(core.PlatformBytedance),
		core.WithHostIP("10.0.0.8"),
		core.WithSender(core.SenderFunc(func(_ context.Context, e core.Event) error {
			notified = true
			if e.Code != 40102 {
				t.Fatalf("unexpected code %d", e.Code)
			}
			if e.HostIP != "10.0.0.8" {
				t.Fatalf("unexpected host ip %q", e.HostIP)
			}
			close(done)
			return nil
		})),
	)

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"?access_token=abc", nil)
	resp, err := tp.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte(body)) {
		t.Fatalf("body mismatch: %s", string(got))
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("expected notify")
	}
	if !notified {
		t.Fatal("expected notify flag")
	}
}

func TestTransportDisabledNoPeekSideEffect(t *testing.T) {
	body := `{"code":40102,"message":"token expired"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, body)
	}))
	defer srv.Close()

	tp := core.NewTransport(http.DefaultTransport, core.WithEnabled(false))
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	resp, err := tp.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

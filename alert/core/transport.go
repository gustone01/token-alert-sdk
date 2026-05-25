package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Transport struct {
	base    http.RoundTripper
	handler *Handler
}

func NewTransport(base http.RoundTripper, opts ...Option) *Transport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &Transport{
		base:    base,
		handler: NewHandler(opts...),
	}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil || resp == nil || resp.Body == nil || t.handler == nil || !t.handler.Enabled() {
		return resp, err
	}

	peek := t.handler.opts.PeekBytes
	if peek <= 0 {
		peek = DefaultPeekBytes
	}

	limited := io.LimitReader(resp.Body, int64(peek))
	buf, _ := io.ReadAll(limited)
	resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(buf), resp.Body))

	var jr struct {
		Code      int64  `json:"code"`
		Message   string `json:"message"`
		MessageCn string `json:"message_cn"`
	}
	if json.Unmarshal(buf, &jr) != nil || jr.Code == 0 {
		return resp, err
	}

	msg := jr.MessageCn
	if msg == "" {
		msg = jr.Message
	}
	tokenHint := req.URL.Query().Get("access_token")
	t.handler.Handle(t.handler.opts.Platform, jr.Code, msg, req.URL.Path, tokenHint)
	return resp, err
}

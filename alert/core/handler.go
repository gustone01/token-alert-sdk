package core

import "context"

type Sender interface {
	Send(ctx context.Context, event Event) error
}

type Handler struct {
	opts Options
}

func NewHandler(opts ...Option) *Handler {
	o := applyOptions(opts)
	return &Handler{opts: o}
}

func (h *Handler) Enabled() bool {
	return h != nil && h.opts.Enabled && h.platformEnabled(h.opts.Platform)
}

func (h *Handler) platformEnabled(p Platform) bool {
	if h.opts.PlatformEnabledFn != nil {
		return h.opts.PlatformEnabledFn(p)
	}
	return h.opts.Enabled
}

func (h *Handler) Handle(platform Platform, code int64, message, apiPath, tokenHint string) {
	if h == nil || !h.opts.Enabled || !h.platformEnabled(platform) {
		return
	}
	if !IsTokenError(platform, code) {
		return
	}
	if h.opts.Sender == nil {
		return
	}

	event := Event{
		Platform:   platform,
		TokenType:  TokenTypeAccessToken,
		Service:    h.opts.ServiceName,
		APIPath:    apiPath,
		Code:       code,
		Message:    message,
		TokenHint:  tokenHint,
		OccurredAt: timeNow(),
	}
	if h.opts.EnrichFunc != nil && tokenHint != "" {
		event.Enrich = h.opts.EnrichFunc(tokenHint)
	}

	go h.notifyAsync(event)
}

func (h *Handler) notifyAsync(event Event) {
	if h.shouldDedup(event) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), h.opts.NotifyTimeout)
	defer cancel()
	_ = h.opts.Sender.Send(ctx, event)
}

func (h *Handler) shouldDedup(event Event) bool {
	if h.opts.Redis == nil {
		return false
	}
	key := dedupKey(event)
	ok, err := h.opts.Redis.SetNX(context.Background(), key, "1", h.opts.DedupTTL).Result()
	if err != nil {
		return false
	}
	return !ok
}

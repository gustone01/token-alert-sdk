package core

import "sync"

var (
	handlers   = make(map[Platform]*Handler)
	handlersMu sync.RWMutex
)

func RegisterHandler(platform Platform, h *Handler) {
	handlersMu.Lock()
	handlers[platform] = h
	handlersMu.Unlock()
}

func HandlerFor(platform Platform) *Handler {
	handlersMu.RLock()
	defer handlersMu.RUnlock()
	return handlers[platform]
}

func TryNotify(platform Platform, code int64, message, apiPath, tokenHint string) {
	h := HandlerFor(platform)
	if h == nil {
		return
	}
	h.Handle(platform, code, message, apiPath, tokenHint)
}

func TryNotifyEvent(event Event) {
	if event.TokenType == "" {
		event.TokenType = TokenTypeAccessToken
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = timeNow()
	}
	h := HandlerFor(event.Platform)
	if h == nil {
		return
	}
	if !h.opts.Enabled || !h.platformEnabled(event.Platform) || !IsTokenError(event.Platform, event.Code) || h.opts.Sender == nil {
		return
	}
	go h.notifyAsync(event)
}

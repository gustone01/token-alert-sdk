package core

import "context"

type SenderFunc func(ctx context.Context, event Event) error

func (f SenderFunc) Send(ctx context.Context, event Event) error {
	return f(ctx, event)
}

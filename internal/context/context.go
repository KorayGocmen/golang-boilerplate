package context

import "context"

type Ctx context.Context

func Background() Ctx {
	return context.Background()
}

func WithValue(parent Ctx, key, val interface{}) Ctx {
	return context.WithValue(parent, key, val)
}

func WithCancel(parent Ctx) (Ctx, context.CancelFunc) {
	return context.WithCancel(parent)
}

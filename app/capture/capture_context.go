package capture

import (
	"context"
	"log"
)

type ctxKey string

const (
	captureKey ctxKey = "capture"
)

// contextManager handle user through the context
type contextManager struct{}

// WithCapture will return a new context with the capture value added to it.
func (c *contextManager) WithCapture(ctx context.Context, capt *Capture) context.Context {
	return context.WithValue(ctx, captureKey, capt)
}

// Get will return the capture assigned to the context, or nil if there
// is any error or there isn't a user.
func (c *contextManager) Get(ctx context.Context) *Capture {
	tmp := ctx.Value(captureKey)
	if tmp == nil {
		return nil
	}
	capt, ok := tmp.(*Capture)
	if !ok {
		log.Printf("context: capture value set incorrectly. type=%T, value=%#v", tmp, tmp)
		return nil
	}
	return capt
}

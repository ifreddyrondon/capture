package capture

import (
	"context"
	"errors"
)

type ctxKey string

const (
	captureKey ctxKey = "capture"
)

var (
	errMissingCapture    = errors.New("capture not found in context")
	errWrongCaptureValue = errors.New("capture value set incorrectly in context")
)

// ContextManager handle capture through the context
type ContextManager struct{}

// NewContextManager returns a new instance of ContextManager
func NewContextManager() *ContextManager {
	return &ContextManager{}
}

// WithCapture will return a new context with the capture value added to it.
func (c *ContextManager) WithCapture(ctx context.Context, capt *Capture) context.Context {
	return context.WithValue(ctx, captureKey, capt)
}

// GetCapture will return the capture assigned to the context, or nil if there
// is any error or there isn't a user.
func (c *ContextManager) GetCapture(ctx context.Context) (*Capture, error) {
	tmp := ctx.Value(captureKey)
	if tmp == nil {
		return nil, errMissingCapture
	}
	capt, ok := tmp.(*Capture)
	if !ok {
		return nil, errWrongCaptureValue
	}
	return capt, nil
}

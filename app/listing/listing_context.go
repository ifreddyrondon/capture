package listing

import (
	"context"
	"errors"
)

type ctxKey string

const (
	paramsKey ctxKey = "listing_params"
)

var (
	errMissingListingParams    = errors.New("listing params not found in context")
	errWrongListingValueParams = errors.New("listing params value set incorrectly in context")
)

// ContextManager handle the listing object through the context
type ContextManager struct{}

// NewContextManager returns a new instance of ContextManager
func NewContextManager() *ContextManager {
	return &ContextManager{}
}

func (c *ContextManager) withParams(ctx context.Context, p *Params) context.Context {
	return context.WithValue(ctx, paramsKey, p)
}

// GetParams will return the listing params reference assigned to the context, or nil if there
// is any error or there isn't a Params.
func (c *ContextManager) GetParams(ctx context.Context) (*Params, error) {
	tmp := ctx.Value(paramsKey)
	if tmp == nil {
		return nil, errMissingListingParams
	}
	p, ok := tmp.(*Params)
	if !ok {
		return nil, errWrongListingValueParams
	}
	return p, nil
}

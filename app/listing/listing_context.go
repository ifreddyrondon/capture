package listing

import (
	"context"
	"errors"
)

type ctxKey string

const (
	paramsKey ctxKey = "listing_value"
)

var (
	errMissingListing    = errors.New("listing not found in context")
	errWrongListingValue = errors.New("listing value set incorrectly in context")
)

// ContextManager handle the listing object through the context
type ContextManager struct{}

// NewContextManager returns a new instance of ContextManager
func NewContextManager() *ContextManager {
	return &ContextManager{}
}

func (c *ContextManager) withParams(ctx context.Context, l *Listing) context.Context {
	return context.WithValue(ctx, paramsKey, l)
}

// GetListing will return the listing reference assigned to the context, or nil if there
// is any error or there isn't a Listing instance.
func (c *ContextManager) GetListing(ctx context.Context) (*Listing, error) {
	tmp := ctx.Value(paramsKey)
	if tmp == nil {
		return nil, errMissingListing
	}
	l, ok := tmp.(*Listing)
	if !ok {
		return nil, errWrongListingValue
	}
	return l, nil
}

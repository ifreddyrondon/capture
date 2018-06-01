package user

import (
	"context"
	"log"
)

type ctxKey string

const (
	userKey ctxKey = "user"
)

// ContextManager handle user through the context
type ContextManager struct{}

// NewContextManager returns a new instance of ContextManager
func NewContextManager() *ContextManager {
	return &ContextManager{}
}

// WithUser will return a new context with the user value added to it.
func (c *ContextManager) WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// Get will return the user assigned to the context, or nil if there
// is any error or there isn't a user.
func (c *ContextManager) Get(ctx context.Context) *User {
	tmp := ctx.Value(userKey)
	if tmp == nil {
		return nil
	}
	user, ok := tmp.(*User)
	if !ok {
		log.Printf("context: user value set incorrectly. type=%T, value=%#v", tmp, tmp)
		return nil
	}
	return user
}

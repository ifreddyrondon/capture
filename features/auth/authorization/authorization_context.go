package authorization

import (
	"context"
	"log"
)

type ctxKey string

const (
	authSubjectIDKey ctxKey = "auth_subject_id_ctx_key"
)

// ContextManager handle user through the context
type ContextManager struct{}

// NewContextManager returns a new instance of ContextManager
func NewContextManager() *ContextManager {
	return &ContextManager{}
}

// WithSubjectID will return a new context with the subject id value added to it.
func (c *ContextManager) WithSubjectID(ctx context.Context, subjectID string) context.Context {
	return context.WithValue(ctx, authSubjectIDKey, subjectID)
}

// Get will return the subject id assigned to the context, or nil if there
// is any error or there isn't a subject id.
func (c *ContextManager) Get(ctx context.Context) string {
	tmp := ctx.Value(authSubjectIDKey)
	if tmp == nil {
		return ""
	}
	subjectID, ok := tmp.(string)
	if !ok {
		log.Printf("context: subject id value set incorrectly. type=%T, value=%#v", tmp, tmp)
		return ""
	}
	return subjectID
}

package authorization

import (
	"context"
	"log"
)

type ctxKey string

const (
	authSubjectIDKey ctxKey = "auth_subject_id_ctx_key"
)

// WithSubjectID will return a new context with the subject id value added to it.
func WithSubjectID(ctx context.Context, subjectID string) context.Context {
	return context.WithValue(ctx, authSubjectIDKey, subjectID)
}

// GetSubjectID will return the subject id assigned to the context, or nil if there
// is any error or there isn't a subject id.
func GetSubjectID(ctx context.Context) string {
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

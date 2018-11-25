package authorization

import (
	"context"
	"errors"
)

type ctxKey string

const (
	authSubjectIDKey ctxKey = "auth_subject_id_ctx_key"
)

var (
	errMissingSubjectID    = errors.New("subject id not found in context")
	errWrongSubjectIDValue = errors.New("subject id value set incorrectly in context")
)

// WithSubjectID will return a new context with the subject id value added to it.
func WithSubjectID(ctx context.Context, subjectID string) context.Context {
	return context.WithValue(ctx, authSubjectIDKey, subjectID)
}

// GetSubjectID will return the subject id assigned to the context, or nil if there
// is any error or there isn't a subject id.
func GetSubjectID(ctx context.Context) (string, error) {
	tmp := ctx.Value(authSubjectIDKey)
	if tmp == nil {
		return "", errMissingSubjectID
	}
	subjectID, ok := tmp.(string)
	if !ok {
		return "", errWrongSubjectIDValue
	}
	return subjectID, nil
}

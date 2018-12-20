package middleware_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/stretchr/testify/assert"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

type notFound string

func (u notFound) Error() string  { return fmt.Sprintf(string(u)) }
func (u notFound) NotFound() bool { return true }

type invalidErr string

func (i invalidErr) Error() string   { return fmt.Sprintf(string(i)) }
func (i invalidErr) IsInvalid() bool { return true }

type notAllowedErr string

func (i notAllowedErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i notAllowedErr) IsNotAuthorized() bool { return true }

func TestContextKeyString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "capture/middleware context value Repository", middleware.RepoCtxKey.String())
}

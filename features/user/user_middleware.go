package user

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/features/auth/authorization"
	"gopkg.in/src-d/go-kallax.v1"
)

// LoggedUser save the authenticated user in a request context.
func LoggedUser(service GetterService) func(next http.Handler) http.Handler {
	json := render.NewJSON()
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			subjectID, err := authorization.GetSubjectID(r.Context())
			if err != nil {
				json.InternalServerError(w, err)
				return
			}
			userID, err := kallax.NewULIDFromText(subjectID)
			if err != nil {
				json.BadRequest(w, err)
				return
			}
			u, err := service.GetByID(userID)
			if err != nil {
				if err == ErrNotFound {
					json.NotFound(w, err)
					return
				}
				json.InternalServerError(w, err)
				return
			}
			ctx := WithUser(r.Context(), u)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

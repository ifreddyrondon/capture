package authorization

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
)

// IsAuthorizedREQ validates if a request can access the resource
func IsAuthorizedREQ(strategy Strategy) func(next http.Handler) http.Handler {
	json := render.NewJSON()
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			subjectID, err := strategy.IsAuthorizedREQ(r)
			if err != nil {
				if strategy.IsNotAuthorizedErr(err) {
					httpErr := render.HTTPError{
						Status:  http.StatusForbidden,
						Error:   http.StatusText(http.StatusForbidden),
						Message: err.Error(),
					}
					json.Response(w, http.StatusForbidden, httpErr)
					return
				}
				json.InternalServerError(w, err)
				return
			}
			ctx := WithSubjectID(r.Context(), subjectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

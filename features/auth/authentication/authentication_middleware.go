package authentication

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/features/user"
)

// Authenticate validate if an user is authorized to continue or 401.
func Authenticate(strategy Strategy) func(next http.Handler) http.Handler {
	json := render.NewJSON()
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u, err := strategy.Validate(r)
			if err != nil {
				if strategy.IsErrCredentials(err) {
					httpErr := render.HTTPError{
						Status:  http.StatusUnauthorized,
						Error:   http.StatusText(http.StatusUnauthorized),
						Message: err.Error(),
					}
					json.Response(w, http.StatusUnauthorized, httpErr)
					return
				}
				if strategy.IsErrDecoding(err) {
					json.BadRequest(w, err)
					return
				}
				json.InternalServerError(w, err)
				return
			}
			ctx := user.WithUser(r.Context(), u)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

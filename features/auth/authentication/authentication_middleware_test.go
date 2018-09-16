package authentication_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/pkg/errors"

	"github.com/go-chi/chi"

	"github.com/ifreddyrondon/capture/features/auth/authentication"
	"github.com/ifreddyrondon/capture/features/user"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
})

func setupStrategy(mock authentication.Strategy) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(authentication.Authenticate(mock))
		r.Post("/", handler)
	})

	return app
}

type strategy struct {
	usr           *user.User
	err           error
	credentialErr bool
	decodingErr   bool
}

func (s *strategy) Validate(r *http.Request) (*user.User, error) { return s.usr, s.err }
func (s *strategy) IsErrCredentials(err error) bool              { return s.credentialErr }
func (s *strategy) IsErrDecoding(err error) bool                 { return s.decodingErr }

func TestAuthenticateSuccess(t *testing.T) {
	t.Parallel()

	strategy := &strategy{usr: &user.User{}}
	app := setupStrategy(strategy)
	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"email": "bla@example.com", "password": "123"}
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		Text().Equal("OK")
}

func TestTokenAuthFailure(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		strategy authentication.Strategy
		payload  map[string]interface{}
		status   int
		response map[string]interface{}
	}{
		{
			name:     "Unauthorized",
			strategy: &strategy{err: errors.New("invalid email or password"), credentialErr: true},
			payload:  map[string]interface{}{"email": "bla@example.com", "password": "123"},
			status:   http.StatusUnauthorized,
			response: map[string]interface{}{
				"status":  401.0,
				"error":   "Unauthorized",
				"message": "invalid email or password",
			},
		},
		{
			name:     "Bad Request",
			strategy: &strategy{err: errors.New("email must not be blank\npassword must not be blank"), decodingErr: true},
			payload:  map[string]interface{}{"email": "bla@example.com", "password": "123"},
			status:   http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "email must not be blank\npassword must not be blank",
			},
		},
		{
			name:     "Internal Server Error",
			strategy: &strategy{err: errors.New("looks like something went wrong")},
			payload:  map[string]interface{}{"email": "bla@example.com", "password": "123"},
			status:   http.StatusInternalServerError,
			response: map[string]interface{}{
				"status":  500.0,
				"error":   "Internal Server Error",
				"message": "looks like something went wrong",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app := setupStrategy(tc.strategy)
			e := bastion.Tester(t, app)
			e.POST("/auth/token-auth").
				WithJSON(tc.payload).
				Expect().
				Status(tc.status).
				JSON().Object().Equal(tc.response)
		})
	}
}

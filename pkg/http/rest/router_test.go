package rest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/http/rest"

	"github.com/ifreddyrondon/bastion"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"

	bastionListing "github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/listing"
	"github.com/ifreddyrondon/capture/pkg/signup"
	"github.com/ifreddyrondon/capture/pkg/updating"

	"github.com/sarulabs/di"
)

type mockSignUpService struct {
	usr *signup.User
	err error
}

func (s *mockSignUpService) EnrollUser(signup.Payload) (*signup.User, error) { return s.usr, s.err }

type mockAuthenticatingService struct {
	usr      *domain.User
	token    string
	err      error
	tokenErr error
}

func (s *mockAuthenticatingService) AuthenticateUser(authenticating.BasicCredential) (*domain.User, error) {
	return s.usr, s.err
}
func (s *mockAuthenticatingService) GetUserToken(kallax.ULID) (string, error) {
	return s.token, s.tokenErr
}

type mockAuthorizingService struct {
	usr *domain.User
	err error
}

func (m *mockAuthorizingService) AuthorizeRequest(*http.Request) (*domain.User, error) {
	return m.usr, m.err
}

type mockRepoService struct {
	repo *domain.Repository
	err  error
}

func (m *mockRepoService) CreateRepo(*domain.User, creating.Payload) (*creating.Repository, error) {
	return &creating.Repository{}, m.err
}
func (m *mockRepoService) GetUserRepos(*domain.User, *bastionListing.Listing) (*listing.ListRepositoryResponse, error) {
	return &listing.ListRepositoryResponse{}, m.err
}

func (m *mockRepoService) GetPublicRepos(*bastionListing.Listing) (*listing.ListRepositoryResponse, error) {
	return &listing.ListRepositoryResponse{}, m.err
}
func (m *mockRepoService) Get(kallax.ULID) (*domain.Repository, error) {
	return m.repo, m.err
}

type mockCaptureService struct {
	capt     *domain.Capture
	captures []domain.Capture
	err      error
}

func (m *mockCaptureService) AddCapture(*domain.Repository, adding.Capture) (*domain.Capture, error) {
	return m.capt, m.err
}
func (m *mockCaptureService) AddCaptures(*domain.Repository, adding.MultiCapture) ([]domain.Capture, error) {
	return m.captures, m.err
}
func (m *mockCaptureService) ListRepoCaptures(*domain.Repository, *bastionListing.Listing) (*listing.ListCaptureResponse, error) {
	return &listing.ListCaptureResponse{}, m.err
}
func (m *mockCaptureService) Get(kallax.ULID, *domain.Repository) (*domain.Capture, error) {
	return m.capt, m.err
}
func (m *mockCaptureService) Update(updating.Capture, *domain.Capture) error { return m.err }
func (m *mockCaptureService) Remove(*domain.Capture) error                   { return m.err }

func resources() di.Container {
	builder, _ := di.NewBuilder()
	definitions := []di.Def{
		{
			Name:  "sign_up-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockSignUpService{}, nil },
		},
		{
			Name:  "authenticating-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockAuthenticatingService{}, nil },
		},
		{
			Name:  "authorize-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockAuthorizingService{}, nil },
		},
		{
			Name:  "creating-repo-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockRepoService{}, nil },
		},
		{
			Name:  "listing-repo-services",
			Build: func(ctn di.Container) (interface{}, error) { return &mockRepoService{}, nil },
		},
		{
			Name:  "getting-repo-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockRepoService{}, nil },
		},
		{
			Name:  "adding-capture-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockCaptureService{}, nil },
		},
		{
			Name:  "adding-multi-capture-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockCaptureService{}, nil },
		},
		{
			Name:  "listing-capture-services",
			Build: func(ctn di.Container) (interface{}, error) { return &mockCaptureService{}, nil },
		},
		{
			Name:  "getting-capture-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockCaptureService{}, nil },
		},
		{
			Name:  "removing-capture-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockCaptureService{}, nil },
		},
		{
			Name:  "updating-capture-service",
			Build: func(ctn di.Container) (interface{}, error) { return &mockCaptureService{}, nil },
		},
	}

	builder.Add(definitions...)
	return builder.Build()
}

func setup() *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Mount("/", rest.Router(resources()))
	return app
}

func TestRouter(t *testing.T) {
	t.Parallel()
	e := bastion.Tester(t, setup())

	tt := []struct {
		uri    string
		method string
	}{
		{uri: "/sign/", method: "POST"},
		{uri: "/auth/token-auth", method: "POST"},
		{uri: "/user/repos/", method: "POST"},
		{uri: "/user/repos/", method: "GET"},
		{uri: "/repositories/", method: "GET"},
		{uri: "/repositories/123", method: "GET"},
		{uri: "/repositories/123/captures", method: "POST"},
		{uri: "/repositories/123/captures/multi", method: "POST"},
		{uri: "/repositories/123/captures", method: "GET"},
		{uri: "/repositories/123/captures/abc", method: "GET"},
		{uri: "/repositories/123/captures/abc", method: "DELETE"},
		{uri: "/repositories/123/captures/abc", method: "PUT"},
	}

	for _, tc := range tt {
		t.Run(tc.uri, func(t *testing.T) {
			r := e.Request(tc.method, tc.uri).Expect().Raw()
			assert.NotEqual(t, 404, r.StatusCode)
			assert.NotEqual(t, 405, r.StatusCode)
		})
	}
}

package repository_test

import (
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/config"
	"github.com/ifreddyrondon/capture/features/repository"
)

func setupRepoController(t *testing.T) (*bastion.Bastion, func()) {
	cfg, err := config.FromString(`PG="postgres://localhost/captures_app_test?sslmode=disable"`)
	if err != nil {
		t.Fatal(err)
	}

	store := cfg.Resources.Get("repository-store").(repository.Store)
	service := repository.Service{Store: store}
	app := bastion.New()
	app.APIRouter.Mount("/repositories/", repository.Routes(service, authOK, loggedUser))

	return app, func() { store.Drop() }
}

func TestListPublicReposWhenEmpty(t *testing.T) {
	app, teardown := setupRepoController(t)
	defer teardown()

	e := bastion.Tester(t, app)
	res := e.GET("/repositories/").Expect().
		JSON().Object()
	res.Value("results").Array().Empty()
	res.Value("listing").Object().
		ContainsKey("paging").
		ContainsKey("sorting")
}

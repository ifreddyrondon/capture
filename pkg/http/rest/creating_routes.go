package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

// SignUp returns a configured http.Handler with creating resources.
func Creating(service creating.Service, auth func(next http.Handler) http.Handler) http.Handler {
	c := &creatingController{service: service, render: render.NewJSON()}

	r := bastion.NewRouter()
	r.Use(auth)
	r.Post("/", c.create)
	return r
}

type creatingController struct {
	service creating.Service
	render  render.APIRenderer
}

func (c *creatingController) create(w http.ResponseWriter, r *http.Request) {
	var payl creating.Payload
	err := creating.Validator.Decode(r, &payl)
	if err != nil {
		c.render.BadRequest(w, err)
		return
	}

	u, err := middleware.GetUser(r.Context())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.render.InternalServerError(w, err)
		return
	}

	repo, err := c.service.CreateRepo(u, payl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.render.InternalServerError(w, err)
		return
	}

	c.render.Created(w, repo)
}

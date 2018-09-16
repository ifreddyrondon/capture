package repository

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/features/user"
)

// Routes returns a configured http.Handler with repository resources.
func Routes(service Service, isAuth, loggedUser func(http.Handler) http.Handler) http.Handler {
	c := &controller{
		service: service,
		render:  render.NewJSON(),
	}

	r := bastion.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(isAuth)
		r.Use(loggedUser)
		r.Post("/", c.create)
	})

	return r
}

// Controller handler the repository routes
type controller struct {
	service Service
	render  render.APIRenderer
}

func (c *controller) create(w http.ResponseWriter, r *http.Request) {
	var repo Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		c.render.BadRequest(w, err)
		return
	}

	userID, err := user.GetUserID(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	repo.UserID = userID

	if err := c.service.Save(&repo); err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	c.render.Created(w, repo)
}

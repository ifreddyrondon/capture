package branch

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

// Controller handler the branch's routes
type Controller struct {
	render render.APIRenderer
}

// NewController returns a new Controller
func NewController() *Controller {
	return &Controller{render: render.NewJSON()}
}

// Router creates a REST router for the branch resource
func (h *Controller) Router() http.Handler {
	r := bastion.NewRouter()
	r.Get("/", h.list)
	return r
}

func (h *Controller) list(w http.ResponseWriter, r *http.Request) {
	b := Branch{}
	h.render.Send(w, b)
}

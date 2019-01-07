package branch

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/domain"
)

// Routes returns a configured http.Handler with branch resources.
func Routes() http.Handler {
	c := &controller{render: render.NewJSON()}

	r := bastion.NewRouter()
	r.Get("/", c.list)
	return r
}

type controller struct {
	render render.APIRenderer
}

func (h *controller) list(w http.ResponseWriter, r *http.Request) {
	b := domain.Branch{}
	h.render.Send(w, b)
}

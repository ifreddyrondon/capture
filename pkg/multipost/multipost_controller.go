package multipost

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/multipost/decoder"
	"github.com/ifreddyrondon/capture/pkg/multipost/encoder"
)

// Routes returns a configured http.Handler with capture resources.
func Routes() http.Handler {
	c := &controller{render: render.NewJSON()}

	r := bastion.NewRouter()
	r.Post("/captures", c.createCaptures)
	return r
}

// Controller handler the capture's routes
type controller struct {
	render render.APIRenderer
}

func (c *controller) createCaptures(w http.ResponseWriter, r *http.Request) {
	var multiPostCaptures decoder.MultiPOSTCaptures
	if err := decoder.Decode(r, &multiPostCaptures); err != nil {
		c.render.BadRequest(w, err)
		return
	}

	// TODO: create job to save captures (async)

	c.render.Send(w, encoder.NewMultiPOSTCaptureResponse(multiPostCaptures))
	return
}

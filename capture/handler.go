package capture

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/utils"
	"github.com/ifreddyrondon/gocapture/app"
	"gopkg.in/mgo.v2"
)

type Handler struct{}

func (h *Handler) Pattern() string {
	return app.CaptureDomain
}

// Routes creates a REST router for the capture resource
func (h *Handler) Router() chi.Router {
	r := gobastion.NewRouter()
	r.Post("/", h.create)
	return r
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	capture := new(Capture)
	if err := utils.ReadJSON(r.Body, capture); err != nil {
		utils.BadRequest(w, err)
		return
	}

	ctx := r.Context()
	err := capture.create(ctx.Value("DB").(*mgo.Database))
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	utils.Created(w, capture)
}

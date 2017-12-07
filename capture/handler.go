package capture

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/utils"
)

type Handler struct{}

// Routes creates a REST router for the capture resource
func (h *Handler) Routes() chi.Router {
	r := gobastion.NewRouter()

	r.Post("/", h.Create)

	return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	path := new(Path)
	if err := utils.ReadJSON(r.Body, path); err != nil {
		utils.BadRequest(w, err)
	}
	utils.Send(w, path)
}

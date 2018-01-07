package branch

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/utils"
	"github.com/ifreddyrondon/gocapture/app"
)

type Handler struct{}

func (h *Handler) Pattern() string {
	return app.BranchDomain
}

// Routes creates a REST router for the branch resource
func (h *Handler) Router() chi.Router {
	r := gobastion.NewRouter()
	r.Post("/", h.Create)
	return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	path := new(Branch)
	if err := utils.ReadJSON(r.Body, path); err != nil {
		utils.BadRequest(w, err)
	}
	utils.Send(w, path)
}

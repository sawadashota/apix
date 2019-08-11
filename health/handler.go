package health

import (
	"net/http"

	"github.com/sawadashota/apix"
)

const (
	// AliveCheckPath is the path where information about the life state of the instance is provided.
	AliveCheckPath = "/health/alive"
)

type healthStatus struct {
	// Status always contains "ok".
	Status string `json:"status"`
}

// Handler handles HTTP requests to health and version endpoints.
type Handler struct {
	srv apix.Server
}

// NewHandler instantiates a handler.
func NewHandler(srv *apix.Server) *Handler {
	return &Handler{
		srv: *srv,
	}
}

// Alive returns an ok status if the instance is ready to handle HTTP requests.
func (h *Handler) Alive(w http.ResponseWriter, r *http.Request) {
	h.srv.Writer().Write(w, r, &healthStatus{
		Status: "ok",
	})
}

// SetRoutes registers this handler's routes.
func (h *Handler) SetRoutes(r *apix.Router) {
	r.GET(AliveCheckPath, h.Alive)
}

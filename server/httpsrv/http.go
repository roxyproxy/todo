package httpsrv

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"todo/config"
	"todo/logger"
	"todo/model"
	"todo/service"

	"github.com/go-chi/chi"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server represents server struct.
type Server struct {
	Serve   http.Handler
	config  *config.Config
	log     logger.Logger
	service service.Handlers
}

// NewHTTPServer returns htts server object.
func NewHTTPServer(service service.Handlers, c *config.Config, l logger.Logger) *Server {
	t := new(Server)
	// t.storage = storage
	t.config = c
	t.log = l
	t.service = service

	s := chi.NewRouter()
	s.Get("/todos", Chain(t.getAllItemsHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Post("/todos", Chain(t.addItemHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Get("/todos/{todoId}", Chain(t.getItemHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Delete("/todos/{todoId}", Chain(t.deleteItemHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Put("/todos/{todoId}", Chain(t.updateItemHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Get("/users", Chain(t.getAllUsersHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Post("/users", Chain(t.addUserHandler, t.SetContentType(), t.Log()))
	s.Get("/users/{userId}", Chain(t.getUserHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Delete("/users/{userId}", Chain(t.deleteUserHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Put("/users/{userId}", Chain(t.updateUserHandler, t.SetContentType(), t.Authorize(), t.Log()))
	s.Post("/user/login", Chain(t.loginUserHandler, t.SetContentType(), t.Log()))

	s.Handle("/metrics", promhttp.Handler())

	t.Serve = s
	return t
}

func (t *Server) handleError(err error, w http.ResponseWriter) {
	status := http.StatusBadRequest
	var statusText string

	t.log.Errorf("HTTP: %s", err)

	switch {
	case errors.Is(err, model.ErrOperational):
		status = http.StatusInternalServerError
	case errors.Is(err, model.ErrBadRequest):
		status = http.StatusBadRequest
	case errors.Is(err, model.ErrNotFound):
		status = http.StatusForbidden
	case errors.Is(err, model.ErrUnauthorized):
		status = http.StatusUnauthorized
	default:
		status = http.StatusBadRequest
	}

	statusText = http.StatusText(status)

	http.Error(w, statusText, status)
}

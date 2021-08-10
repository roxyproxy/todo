package server

import (
	"net/http"
	"todo/config"
	"todo/storage"

	"github.com/go-chi/chi"
)

type TodoServer struct {
	storage storage.Storage
	Serve   http.Handler
	config  config.Config
}

func NewTodoServer(storage storage.Storage, c config.Config) *TodoServer {
	t := new(TodoServer)
	t.storage = storage
	t.config = c
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

	t.Serve = s
	return t
}

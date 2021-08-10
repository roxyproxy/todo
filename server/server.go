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

	s.Get("/todos", Chain(t.getAllItemsHandler, SetContentType(), Authorize(), Log()))
	s.Post("/todos", Chain(t.addItemHandler, SetContentType(), Authorize(), Log()))
	s.Get("/todos/{todoId}", Chain(t.getItemHandler, SetContentType(), Authorize(), Log()))
	s.Delete("/todos/{todoId}", Chain(t.deleteItemHandler, SetContentType(), Authorize(), Log()))
	s.Put("/todos/{todoId}", Chain(t.updateItemHandler, SetContentType(), Authorize(), Log()))

	s.Get("/users", Chain(t.getAllUsersHandler, SetContentType(), Authorize(), Log()))
	s.Post("/users", Chain(t.addUserHandler, SetContentType(), Log()))
	s.Get("/users/{userId}", Chain(t.getUserHandler, SetContentType(), Authorize(), Log()))
	s.Delete("/users/{userId}", Chain(t.deleteUserHandler, SetContentType(), Authorize(), Log()))
	s.Put("/users/{userId}", Chain(t.updateUserHandler, SetContentType(), Authorize(), Log()))

	s.Post("/user/login", Chain(t.loginUserHandler, SetContentType(), Log()))

	t.Serve = s
	return t
}

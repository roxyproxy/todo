package main

import (
	"net/http"
	"todo/server"
	"todo/storage/inmemory"
)

func main() {
	server := server.NewTodoServer(inmemory.NewInMemoryStorage())

	http.ListenAndServe(":5008", server)
}

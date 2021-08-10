package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	conf "todo/config"
	"todo/server"
	"todo/storage/inmemory"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	var config = conf.New()

	server := server.NewTodoServer(inmemory.NewInMemoryStorage(), *config)

	http.ListenAndServe(":5001", server.Serve)
}

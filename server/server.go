package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"todo/model"
	"todo/storage"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const secretKey = "verysecretkey"

type TodoServer struct {
	storage storage.Storage
	Serve   http.Handler
}

func NewTodoServer(storage storage.Storage) *TodoServer {
	t := new(TodoServer)
	t.storage = storage
	s := chi.NewRouter()

	s.Get("/todos", t.getAllItemsHandler)
	s.Post("/todos", t.addItemHandler)
	s.Get("/todos/{todoId}", t.getItemHandler)
	s.Delete("/todos/{todoId}", t.deleteItemHandler)
	s.Put("/todos/{todoId}", t.updateItemHandler)

	s.Get("/users", t.getAllUsersHandler)
	s.Post("/users", t.addUserHandler)
	s.Get("/users/{userId}", t.getUserHandler)
	s.Delete("/users/{userId}", t.deleteUserHandler)
	s.Put("/users/{userId}", t.updateUserHandler)

	s.Post("/user/login", t.loginUserHandler)

	t.Serve = s
	return t
}

// todos handlers
func (t *TodoServer) getAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	filter := storage.TodoFilter{}
	if val, ok := r.URL.Query()["status"]; ok {
		filter.Status = val[0]
	}
	if date, ok := r.URL.Query()["fromdate"]; ok {
		fromDate, err := time.Parse(time.RFC3339, date[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		filter.FromDate = &fromDate
	}
	if date, ok := r.URL.Query()["todate"]; ok {
		toDate, err := time.Parse(time.RFC3339, date[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		filter.FromDate = &toDate
	}

	items, err := t.storage.GetAllItems(filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(items)

}

func (t *TodoServer) addItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todo := model.TodoItem{}
	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := t.storage.AddItem(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j := model.TodoId{Id: id}
	json.NewEncoder(w).Encode(j)
}

func (t *TodoServer) getItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "todoId")
	todo, err := t.storage.GetItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if todo.Id == "" {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

func (t *TodoServer) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "todoId")
	todo, err := t.storage.GetItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if todo.Id == "" {
		http.NotFound(w, r)
		return
	}
	err = t.storage.DeleteItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func (t *TodoServer) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "todoId")
	todo, err := t.storage.GetItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if todo.Id == "" {
		http.NotFound(w, r)
		return
	}
	todo = model.TodoItem{}
	json.NewDecoder(r.Body).Decode(&todo)
	todo.Id = id
	err = t.storage.UpdateItem(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// users handlers
func (t *TodoServer) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := t.storage.GetAllUsers(storage.UserFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (t *TodoServer) addUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password, err = hashPassword(user.Password)

	id, err := t.storage.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j := model.TodoId{Id: id}
	json.NewEncoder(w).Encode(j)
}

func (t *TodoServer) getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "userId")
	user, err := t.storage.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Id == "" {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (t *TodoServer) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "userId")
	user, err := t.storage.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Id == "" {
		http.NotFound(w, r)
		return
	}
	err = t.storage.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func (t *TodoServer) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "userId")
	user, err := t.storage.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Id == "" {
		http.NotFound(w, r)
		return
	}
	user = model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	user.Id = id
	err = t.storage.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (t *TodoServer) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	credentials := model.Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)

	userId, err := t.authenticateUser(credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := generateToken(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(token)

	fmt.Println(t.storage)
}

func (t *TodoServer) authenticateUser(credentials model.Credentials) (string, error) {
	filter := storage.UserFilter{}
	filter.UserName = credentials.UserName

	users, _ := t.storage.GetAllUsers(filter)

	if len(users) == 0 || checkPasswordHash(credentials.Password, users[0].Password) {
		return "", errors.New("Invalid user credentials")
	}
	return users[0].Id, nil
}

func generateToken(id string) (tokenString string, err error) {
	claims := jwt.MapClaims{}
	claims["userid"] = id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
	"todo/model"
	"todo/storage"
)

type TodoServer struct {
	storage storage.Storage
	routes  []model.Route
	http.Handler
}

func NewTodoServer(storage storage.Storage) *TodoServer {
	t := new(TodoServer)
	t.storage = storage
	t.routes = []model.Route{
		{"GET", "/todo", t.getAllItemsHandler},
		//{"GET", "/todo?([^/]+)$", t.getFilteredItemsHandler},
		{"POST", "/todo", t.addItemHandler},
		{"GET", "/todo/([^/]+)", t.getItemHandler},
		{"DELETE", "/todo/([^/]+)", t.deleteItemHandler},
		{"PUT", "/todo/([^/]+)", t.updateItemHandler},
	}
	return t
}

func (t *TodoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range t.routes {
		re := regexp.MustCompile("^" + route.Path + "$")
		matches := re.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 && route.Method == r.Method {
			if len(matches) == 2 {
				q := r.URL.Query()
				q.Add("id", matches[1])
				r.URL.RawQuery = q.Encode()
			}
			w.Header().Set("Content-Type", "application/json")
			route.Handlr(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// handlers
func (t *TodoServer) getAllItemsHandler(w http.ResponseWriter, r *http.Request) {
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

func (t *TodoServer) getFilteredItemsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getFilteredItemsHandler")
	/*status := r.URL.Query()["Status"][0]
	fromDate, err := time.Parse(time.RFC3339, r.URL.Query()["fromDate"][0])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	toDate, err := time.Parse(time.RFC3339, r.URL.Query()["toDate"][0])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	*/
	filter := storage.TodoFilter{Status: "Done"}
	items, err := t.storage.GetAllItems(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (t *TodoServer) addItemHandler(w http.ResponseWriter, r *http.Request) {
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
	id := r.URL.Query()["id"][0]
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
	id := r.URL.Query()["id"][0]
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
	id := r.URL.Query()["id"][0]
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

//matching
func match(path string, pattern string) bool {
	re := regexp.MustCompile("^" + pattern + "$")
	matches := re.FindStringSubmatch(path)
	if len(matches) > 0 {
		return true
	}
	return false
}
func getMatchParams(path string, pattern string) []string {
	re := regexp.MustCompile("^" + pattern + "$")
	matches := re.FindStringSubmatch(path)
	if len(matches) > 0 {
		return matches[1:]
	}
	return []string{}
}

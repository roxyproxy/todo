package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"todo/model"
	"todo/storage"
)

type TodoServer struct {
	storage storage.Storage
	routes  []model.Route
	http.Handler
}

func NewTodoServer(storage storage.Storage) *TodoServer {
	p := new(TodoServer)
	p.storage = storage
	p.routes = []model.Route{
		{"GET", "/todo", p.getAllItemsHandler},
		{"POST", "/todo", p.addItemHandler},
	}

	return p
}

type ctxKey struct{}

func (t *TodoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	p := r.URL.Path
	switch {
	case match(p, "/todo") && r.Method == "GET":
		todos, _ := t.storage.GetAllItems(storage.TodoFilter{})
		fmt.Fprint(w, todos)
	case match(p, "/todo") && r.Method == "POST":
		todo := model.TodoItem{}
		json.NewDecoder(r.Body).Decode(&todo)
		t.storage.AddItem(todo)
		fmt.Fprint(w, todo)
	case match(p, "/todo/([^/]+)") && r.Method == "DELETE":
		params := getMatchParams(p, "/todo/([^/]+)")
		err := t.storage.DeleteItem(params[0])
		fmt.Fprint(w, err)
	case match(p, "/todo/([^/]+)") && r.Method == "GET":
		params := getMatchParams(p, "/todo/([^/]+)")
		todo, _ := t.storage.GetItem(params[0])
		fmt.Fprint(w, todo)
	case match(p, "/todo/([^/]+)") && r.Method == "PUT":
		todo := model.TodoItem{}
		json.NewDecoder(r.Body).Decode(&todo)
		params := getMatchParams(p, "/todo/([^/]+)")
		todo.Id = params[0]
		err := t.storage.UpdateItem(todo)
		fmt.Fprint(w, err)
	default:
		http.NotFound(w, r)
		return
	}

	/*for _, route := range t.routes {
		re := regexp.MustCompile("^" + route.Path + "$")
		matches := re.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			//ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.Handler(w, r)
			return
		}
	}
	http.NotFound(w, r)
	*/

}

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

func (t *TodoServer) testHandller(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "test")
}

func (t *TodoServer) getAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	items, _ := t.storage.GetAllItems(storage.TodoFilter{})
	fmt.Fprint(w, items)
}
func (t *TodoServer) addItemHandler(w http.ResponseWriter, r *http.Request) {
	todo := model.TodoItem{}
	json.NewDecoder(r.Body).Decode(&todo)
	t.storage.AddItem(todo)
	fmt.Fprint(w, todo)
}

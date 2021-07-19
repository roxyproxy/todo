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
	http.Handler
}

func NewTodoServer(storage storage.Storage) *TodoServer {
	p := new(TodoServer)
	p.storage = storage
	//p.Handler = http.HandlerFunc(Router)

	return p
}

func match(path string, re string) bool {
	r := regexp.MustCompile("^" + re + "$")
	return r.MatchString(path)
}

func (t *TodoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	p := r.URL.Path
	switch {
	case match(p, "/todo") && r.Method == "GET":
		items, _ := t.storage.GetAllItems(storage.TodoFilter{})
		fmt.Fprint(w, items)
	case match(p, "/todo") && r.Method == "POST":
		todo := model.TodoItem{}
		json.NewDecoder(r.Body).Decode(&todo)
		t.storage.AddItem(todo)
		fmt.Fprint(w, todo)
	/*case match(p, "/todo/([^/]+)") && r.Method == "GET":
	item, _ := t.storage.GetItem()
	fmt.Fprint(w, item)*/
	default:
		http.NotFound(w, r)
		return
	}

}

/**type route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.HandlerFunc
}

var routes = []route{
	{"GET", "/todo", contact},
	{"POST", "/todo", apiGetWidgets},
}

/*
func Router(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var mutex sync.Mutex

	mutex.Lock()
	defer mutex.Unlock()

	p := r.URL.Path
	switch {
	case match(p, "/todo") && r.Method == "GET":
		//items, _ := t.storage.GetAllItems(storage.TodoFilter{})
		fmt.Fprint(w, "items")
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}
*/

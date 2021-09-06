package httpsrv

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"time"
	"todo/model"
	"todo/storage"
)

// todos handlers.
func (t *Server) getAllItemsHandler(w http.ResponseWriter, r *http.Request) {
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
	err = json.NewEncoder(w).Encode(items)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (t *Server) addItemHandler(w http.ResponseWriter, r *http.Request) {
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
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (t *Server) getItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	todo, err := t.storage.GetItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if todo.Id == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (t *Server) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	todo, err := t.storage.GetItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if todo.Id == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	err = t.storage.DeleteItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (t *Server) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	todo, err := t.storage.GetItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if todo.Id == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	todo = model.TodoItem{}
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.Id = id
	err = t.storage.UpdateItem(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

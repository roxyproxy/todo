package server

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"time"
	"todo/model"
	"todo/storage"
)

// todos handlers
func (t *TodoServer) getAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
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
	//w.Header().Set("Content-Type", "application/json")
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
	//w.Header().Set("Content-Type", "application/json")
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
	json.NewEncoder(w).Encode(todo)
}

func (t *TodoServer) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
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
func (t *TodoServer) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
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
	json.NewDecoder(r.Body).Decode(&todo)
	todo.Id = id
	err = t.storage.UpdateItem(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

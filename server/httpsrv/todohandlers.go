package httpsrv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
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
			t.handleError(fmt.Errorf("%q: %q: %w", "Error in getAllItemsHandler.", err, model.ErrBadRequest), w)
			return
		}
		filter.FromDate = &fromDate
	}
	if date, ok := r.URL.Query()["todate"]; ok {
		toDate, err := time.Parse(time.RFC3339, date[0])
		if err != nil {
			t.handleError(fmt.Errorf("%q: %q: %w", "Error in getAllItemsHandler.", err, model.ErrBadRequest), w)
			return
		}
		filter.ToDate = &toDate
	}

	items, err := t.service.GetTodos(r.Context(), filter)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in getAllItemsHandler.", err), w)
		return
	}
	err = json.NewEncoder(w).Encode(items)

	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in getAllItemsHandler.", err, model.ErrBadRequest), w)
		return
	}
}

func (t *Server) addItemHandler(w http.ResponseWriter, r *http.Request) {
	todo := model.TodoItem{}
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in add todo handler", err, model.ErrBadRequest), w)
		return
	}

	id, err := t.service.AddTodo(r.Context(), todo)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in add todo handler", err), w)
		return
	}

	j := model.TodoID{ID: id}
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in add todo handler", err, model.ErrBadRequest), w)
		return
	}
}

func (t *Server) getItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	todo, err := t.service.GetTodo(r.Context(), id)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in getItemsHandler.", err), w)
		return
	}
	if todo.ID == "" {
		t.handleError(fmt.Errorf("%q: %w", "Error in getItemHandler.", model.ErrNotFound), w)
		return
	}
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in getItemsHandler.", err, model.ErrBadRequest), w)
		return
	}
}

func (t *Server) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	err := t.service.DeleteTodo(r.Context(), id)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in deleteItemHandler.", err), w)
		return
	}
}

func (t *Server) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "todoId")
	todo := model.TodoItem{}
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "updateItemHandler.", err, model.ErrBadRequest), w)
		return
	}
	err = t.service.UpdateTodo(r.Context(), id, todo)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in updateItemHandler.", err), w)
		return
	}
}

package httpsrv

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strings"
	"todo/model"
	"todo/storage"
)

// users handlers.
func (t *Server) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	filter := storage.UserFilter{}
	if val, ok := r.URL.Query()["username"]; ok {
		filter.UserName = val[0]
	}

	users, err := t.service.GetUsers(filter, r.Context())
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in getAllUsersHandler.", err, model.ErrBadRequest), w)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in getAllUsersHandler.", model.ErrBadRequest), w)
		return
	}
}

func (t *Server) addUserHandler(w http.ResponseWriter, r *http.Request) {
	newUser := model.NewUser{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in add user handler", model.ErrBadRequest), w)
		return
	}
	user := model.User{
		UserName:  newUser.UserName,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Password:  newUser.Password,
		Location:  newUser.Location,
	}

	id, err := t.service.AddUser(user)
	j := model.TodoId{Id: id}
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in add user handler", model.ErrBadRequest), w)
		return
	}
}

func (t *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// userid := getUserFromContext(r)

	id := chi.URLParam(r, "userId")
	user, err := t.storage.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Id == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (t *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	user, err := t.storage.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Id == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	err = t.storage.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (t *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	user, err := t.storage.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Id == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		// http.NotFound(w, r)
		return
	}
	user = model.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Id = id
	err = t.storage.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (t *Server) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials model.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in loginUsersHandler.", model.ErrUnauthorized), w)
		return
	}
	token, err := t.service.LoginUser(credentials)

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in loginUsersHandler.", model.ErrUnauthorized), w)
		return
	}
}

func getToken(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	strArr := strings.Split(tokenString, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

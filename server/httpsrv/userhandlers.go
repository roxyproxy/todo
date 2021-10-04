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

	users, err := t.service.GetUsers(r.Context(), filter)

	if err != nil {
		t.handleError(fmt.Errorf("%q: %w: ", "Error in getAllUsersHandler.", err), w)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in getAllUsersHandler.", err, model.ErrBadRequest), w)
		return
	}
}

func (t *Server) addUserHandler(w http.ResponseWriter, r *http.Request) {
	newUser := model.NewUser{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in add user handler", err, model.ErrBadRequest), w)
		return
	}
	user := model.User{
		UserName:  newUser.UserName,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Password:  newUser.Password,
		Location:  newUser.Location,
	}

	id, err := t.service.AddUser(r.Context(), user)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in add user handler", err, model.ErrBadRequest), w)
		return
	}
	j := model.TodoId{Id: id}
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in add user handler", err, model.ErrBadRequest), w)
		return
	}
}

func (t *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	user, err := t.service.GetUser(r.Context(), id)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in get user handler", err), w)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in get user handler", err, model.ErrBadRequest), w)
		return
	}
}

func (t *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	err := t.service.DeleteUser(r.Context(), id)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in delete user handler", err), w)
		return
	}
}

func (t *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in update user handler", err, model.ErrBadRequest), w)
		return
	}

	err = t.service.UpdateUser(r.Context(), id, user)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %w", "Error in update user handler", err), w)
		return
	}
}

func (t *Server) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials model.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in loginUsersHandler.", err, model.ErrUnauthorized), w)
		return
	}
	token, err := t.service.LoginUser(r.Context(), credentials)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in loginUsersHandler.", err, model.ErrUnauthorized), w)
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		t.handleError(fmt.Errorf("%q: %q: %w", "Error in loginUsersHandler.", err, model.ErrUnauthorized), w)
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

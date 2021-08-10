package server

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"todo/model"
	"todo/storage"
)

// users handlers
func (t *TodoServer) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	//userid := getUserFromContext(r)
	filter := storage.UserFilter{}
	if val, ok := r.URL.Query()["username"]; ok {
		filter.UserName = val[0]
	}

	users, err := t.storage.GetAllUsers(filter)
	if err != nil {
		//t.handleError(err, w)  do log, do error cast,
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (t *TodoServer) addUserHandler(w http.ResponseWriter, r *http.Request) {
	credentials := model.Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := model.User{UserName: credentials.UserName}
	user.Password, err = hashPassword(credentials.Password)

	id, err := t.storage.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j := model.TodoId{Id: id}
	json.NewEncoder(w).Encode(j)
}

func (t *TodoServer) getUserHandler(w http.ResponseWriter, r *http.Request) {
	//userid := getUserFromContext(r)

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
	json.NewEncoder(w).Encode(user)
}

func (t *TodoServer) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
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
func (t *TodoServer) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	user, err := t.storage.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Id == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		//http.NotFound(w, r)
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
	var credentials model.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)

	userId, err := t.authenticateUser(credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := generateToken(userId, t.config.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(token)
}

func (t *TodoServer) authenticateUser(credentials model.Credentials) (string, error) {
	filter := storage.UserFilter{credentials.UserName}

	users, _ := t.storage.GetAllUsers(filter)

	if len(users) == 0 || !checkPasswordHash(credentials.Password, users[0].Password) {
		return "", errors.New("Invalid user credentials")
	}
	return users[0].Id, nil
}

func generateToken(id string, secretKey string) (token model.Token, err error) {
	claims := model.Claims{
		id,
		jwt.StandardClaims{},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.TokenString, err = t.SignedString([]byte(secretKey))
	if err != nil {
		return token, err
	}
	return token, nil
}

func getToken(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	strArr := strings.Split(tokenString, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*
func getUserFromContext(r *http.Request) string {
	userid := r.Context().Value("userid")

	return userid.(string)
}
*/

package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"todo/config"
	"todo/model"
	"todo/storage"
)

type Handlers interface {
	//GetTodos(w http.ResponseWriter, r *http.Request)
	//AddTodo(w http.ResponseWriter, r *http.Request)
	//GetTodo(w http.ResponseWriter, r *http.Request)
	//DeleteTodo(w http.ResponseWriter, r *http.Request)
	//UpdateTodo(w http.ResponseWriter, r *http.Request)
	AddUser(user model.User) (string, error)
	//DeleteUser(id string) error
	//UpdateUser(user model.User) error
	//GetUser(id string) (model.User, error)
	GetUsers(filter storage.UserFilter, ctx context.Context) ([]model.User, error)
	LoginUser(credentials model.Credentials) (model.Token, error)
	ValidateToken(tokenString string) (*model.Claims, error)
}

type handlersService struct {
	storage storage.Storage
	config  *config.Config
}

func NewService(storage storage.Storage, c *config.Config) Handlers {
	return &handlersService{storage: storage, config: c}
}

func (h *handlersService) GetUsers(filter storage.UserFilter, ctx context.Context) ([]model.User, error) {
	err := checkUserNotEmptyInContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("%q: %q: %w", "Could not get all users.", err, model.ErrUnauthorized)
	}

	users, err := h.storage.GetAllUsers(filter)
	if err != nil {
		return nil, fmt.Errorf("%q: %q: %w", "Could not get all users.", err, model.ErrOperational)
	}

	return users, nil
}

func (h *handlersService) AddUser(user model.User) (string, error) {
	user.Password, _ = hashPassword(user.Password)

	id, err := h.storage.AddUser(user)
	if err != nil {
		return "", fmt.Errorf("%q: %w", "Could not add user", model.ErrBadRequest)
	}
	return id, nil
}

func (h *handlersService) LoginUser(credentials model.Credentials) (model.Token, error) {
	userId, err := h.authenticateUser(credentials)
	if err != nil {
		return model.Token{}, fmt.Errorf("Login error: %v: %w", err, model.ErrUnauthorized)
	}

	token, err := generateToken(userId, h.config.SecretKey)
	if err != nil {
		return model.Token{}, fmt.Errorf("%v: %w", err, model.ErrUnauthorized)
	}

	return token, nil
}
func checkUserNotEmptyInContext(ctx context.Context) error {
	userid := ctx.Value(model.KeyUserId("userid"))

	if userid == "" {
		return fmt.Errorf("%q", "Userid in context is empty or does not exists.")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *handlersService) authenticateUser(credentials model.Credentials) (string, error) {
	filter := storage.UserFilter{UserName: credentials.UserName}
	users, _ := h.storage.GetAllUsers(filter)

	if len(users) == 0 || !checkPasswordHash(credentials.Password, users[0].Password) {
		return "", fmt.Errorf("%q", "Invalid user credentials.")
	}
	return users[0].Id, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(id string, secretKey string) (token model.Token, err error) {
	claims := model.Claims{UserId: id}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.TokenString, err = t.SignedString([]byte(secretKey))
	if err != nil {
		return token, err
	}
	return token, nil
}

func (h *handlersService) ValidateToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.config.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%q: %q: %w", "invalid token: authentication failed.", err, model.ErrUnauthorized)
	}
	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("%q: %q: %w", "invalid token: authentication failed.", err, model.ErrUnauthorized)

	}
	return claims, nil
}

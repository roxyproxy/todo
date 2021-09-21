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
	GetTodos(ctx context.Context, filter storage.TodoFilter) ([]model.TodoItem, error)
	AddTodo(ctx context.Context, todo model.TodoItem) (string, error)
	GetTodo(ctx context.Context, id string) (model.TodoItem, error)
	DeleteTodo(ctx context.Context, id string) error
	UpdateTodo(ctx context.Context, id string, user model.TodoItem) error

	AddUser(ctx context.Context, user model.User) (string, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, user model.User) error
	GetUser(ctx context.Context, id string) (model.User, error)
	GetUsers(ctx context.Context, filter storage.UserFilter) ([]model.User, error)
	LoginUser(ctx context.Context, credentials model.Credentials) (model.Token, error)

	ValidateToken(ctx context.Context, tokenString string) (*model.Claims, error)
	AuthenticateUser(credentials model.Credentials) (string, error)
	GenerateToken(id string, secretKey string) (token model.Token, err error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type handlersService struct {
	storage storage.Storage
	config  *config.Config
}

func NewService(storage storage.Storage, c *config.Config) Handlers {
	return &handlersService{storage: storage, config: c}
}

func (h *handlersService) GetUsers(ctx context.Context, filter storage.UserFilter) ([]model.User, error) {
	_, err := h.getUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%q %q %w", "Could not get all users.", err, model.ErrUnauthorized)
	}

	users, err := h.storage.GetAllUsers(filter)
	if err != nil {
		return nil, fmt.Errorf("%q: %q: %w", "Could not get all users.", err, model.ErrOperational)
	}

	return users, nil
}

func (h *handlersService) GetUser(ctx context.Context, id string) (model.User, error) {
	_, err := h.getUserFromContext(ctx)
	if err != nil {
		return model.User{}, fmt.Errorf("%q: %q: %w", "Could not get user.", err, model.ErrUnauthorized)
	}

	user, err := h.storage.GetUser(id)
	if err != nil {
		return model.User{}, fmt.Errorf("%q: %q: %w", "Could not get user.", err, model.ErrOperational)
	}
	if user.Id == "" {
		return model.User{}, fmt.Errorf("%q: %q: %w", "Could not get user.", err, model.ErrNotFound)
	}
	return user, nil
}

func (h *handlersService) AddUser(ctx context.Context, user model.User) (string, error) {
	user.Password, _ = h.HashPassword(user.Password)

	id, err := h.storage.AddUser(user)
	if err != nil {
		return "", fmt.Errorf("%q: %w", "Could not add user", model.ErrBadRequest)
	}
	return id, nil
}

func (h *handlersService) DeleteUser(ctx context.Context, id string) error {
	_, err := h.getUserFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not delete user.", err, model.ErrUnauthorized)
	}

	user, err := h.storage.GetUser(id)
	if err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not delete user.", err, model.ErrOperational)
	}
	if user.Id == "" {
		return fmt.Errorf("%q: %q: %w", "Could not delete user.", err, model.ErrNotFound)
	}

	if err := h.storage.DeleteUser(id); err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not get user.", err, model.ErrOperational)
	}
	return nil
}

func (h *handlersService) UpdateUser(ctx context.Context, id string, user model.User) error {
	_, err := h.getUserFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not update user.", err, model.ErrUnauthorized)
	}

	u, err := h.storage.GetUser(id)
	if err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not update user.", err, model.ErrOperational)
	}
	if u.Id == "" {
		return fmt.Errorf("%q: %q: %w", "Could not update user.", err, model.ErrNotFound)
	}
	user.Id = id
	err = h.storage.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("%q: %w", "Could not update user", model.ErrBadRequest)
	}
	return nil
}

func (h *handlersService) LoginUser(ctx context.Context, credentials model.Credentials) (model.Token, error) {
	userId, err := h.AuthenticateUser(credentials)
	if err != nil {
		return model.Token{}, fmt.Errorf("Login error: %v: %w", err, model.ErrUnauthorized)
	}

	token, err := h.GenerateToken(userId, h.config.SecretKey)
	if err != nil {
		return model.Token{}, fmt.Errorf("%v: %w", err, model.ErrUnauthorized)
	}

	return token, nil
}
func (h *handlersService) getUserFromContext(ctx context.Context) (string, error) {
	userid := ctx.Value(model.KeyUserId("userid"))
	if userid == nil {
		return "", fmt.Errorf("%q", "Userid in context is not provided.")
	}
	if userid.(string) == "" {
		return "", fmt.Errorf("%q", "Userid in context is empty.")
	}

	user, err := h.storage.GetUser(userid.(string))
	if err != nil {
		return "", fmt.Errorf("%q", err)
	}
	if user.Id == "" {
		return "", fmt.Errorf("%q", "User in context does not exists.")
	}

	return userid.(string), nil
}

func (h *handlersService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *handlersService) AuthenticateUser(credentials model.Credentials) (string, error) {
	filter := storage.UserFilter{UserName: credentials.UserName}
	users, _ := h.storage.GetAllUsers(filter)

	if len(users) == 0 || !h.CheckPasswordHash(credentials.Password, users[0].Password) {
		return "", fmt.Errorf("%q", "Invalid user credentials.")
	}
	return users[0].Id, nil
}

func (h *handlersService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *handlersService) GenerateToken(id string, secretKey string) (token model.Token, err error) {
	claims := model.Claims{UserId: id}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.TokenString, err = t.SignedString([]byte(secretKey))
	if err != nil {
		return token, err
	}
	return token, nil
}

func (h *handlersService) ValidateToken(ctx context.Context, tokenString string) (*model.Claims, error) {
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

func (h *handlersService) AddTodo(ctx context.Context, todo model.TodoItem) (string, error) {
	userid, err := h.getUserFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("%q: %q: %w", "Could not add todo.", err, model.ErrUnauthorized)
	}
	todo.UserId = userid
	id, err := h.storage.AddItem(todo)
	if err != nil {
		return "", fmt.Errorf("%q: %w", "Could not add todo", model.ErrBadRequest)
	}
	return id, nil
}

func (h *handlersService) GetTodos(ctx context.Context, filter storage.TodoFilter) ([]model.TodoItem, error) {
	userid, err := h.getUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%q: %q: %w", "Could not get all todos.", err, model.ErrUnauthorized)
	}
	filter.UserId = userid

	todos, err := h.storage.GetAllItems(filter)
	if err != nil {
		return nil, fmt.Errorf("%q: %q: %w", "Could not get all users.", err, model.ErrOperational)
	}

	return todos, nil
}

func (h *handlersService) GetTodo(ctx context.Context, id string) (model.TodoItem, error) {
	userid, err := h.getUserFromContext(ctx)
	if err != nil {
		return model.TodoItem{}, fmt.Errorf("%q: %q: %w", "Could not get todo.", err, model.ErrUnauthorized)
	}

	todo, err := h.storage.GetItem(id)
	if err != nil {
		return model.TodoItem{}, fmt.Errorf("%q: %q: %w", "Could not get todo.", err, model.ErrOperational)
	}
	if todo.Id == "" || todo.UserId != userid {
		return model.TodoItem{}, fmt.Errorf("%q: %q: %w", "Could not get todo.", err, model.ErrNotFound)
	}

	return todo, nil
}

func (h *handlersService) DeleteTodo(ctx context.Context, id string) error {
	userid, err := h.getUserFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not delete todo.", err, model.ErrUnauthorized)
	}

	todo, err := h.storage.GetItem(id)
	if err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not delete todo.", err, model.ErrOperational)
	}
	if todo.Id == "" || todo.UserId != userid {
		return fmt.Errorf("%q: %q: %w", "Could not delete todo.", err, model.ErrNotFound)
	}

	if err := h.storage.DeleteItem(id); err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not get tod.", err, model.ErrOperational)
	}
	return nil
}

func (h *handlersService) UpdateTodo(ctx context.Context, id string, todo model.TodoItem) error {
	userid, err := h.getUserFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not update todo.", err, model.ErrUnauthorized)
	}

	u, err := h.storage.GetItem(id)
	if err != nil {
		return fmt.Errorf("%q: %q: %w", "Could not update todo.", err, model.ErrOperational)
	}
	if u.Id == "" || u.UserId != userid {
		return fmt.Errorf("%q: %q: %w", "Could not update todo.", err, model.ErrNotFound)
	}

	todo.Id = id
	todo.UserId = userid
	err = h.storage.UpdateItem(todo)
	if err != nil {
		return fmt.Errorf("%q: %w", "Could not update todo", model.ErrBadRequest)
	}
	return nil
}

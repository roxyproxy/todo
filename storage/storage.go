package storage

import (
	"time"
	"todo/model"
)

// Storage represent interface for storage types.
type Storage interface {
	AddItem(item model.TodoItem) (id string, err error)
	DeleteItem(id string) error
	UpdateItem(item model.TodoItem) error
	GetItem(id string) (model.TodoItem, error)
	GetAllItems(filter TodoFilter) ([]model.TodoItem, error)

	AddUser(user model.User) (id string, err error)
	DeleteUser(id string) error
	UpdateUser(user model.User) error
	GetUser(id string) (model.User, error)
	GetAllUsers(filter UserFilter) ([]model.User, error)
}

// TodoFilter represents filter struct for todos.
type TodoFilter struct {
	FromDate *time.Time // nil if empty ?
	ToDate   *time.Time // nil if empty ?
	Status   string
	UserID   string
}

// UserFilter represents filter struct for users.
type UserFilter struct {
	UserName string
}

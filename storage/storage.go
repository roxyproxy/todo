package storage

import (
	"time"
	"todo/model"
)

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
	//AuthenticateUser(credentials model.Credentials) (id string, err error)
}

type TodoFilter struct {
	FromDate *time.Time //nil if empty
	ToDate   *time.Time
	Status   string
}

type UserFilter struct {
	UserName string
}

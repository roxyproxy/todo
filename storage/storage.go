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
}

type TodoFilter struct {
	FromDate *time.Time //nil if empty
	ToDate   *time.Time
	Status   string
}

package storage

import (
	"time"
	"todo/model"
)

type Storage interface {
	AddItem(item model.TodoItem) (id int64)
	DeleteItem(id int64)
	UpdateItem(item model.TodoItem)
	GetItem(id int64) model.TodoItem
	GetAllItems(filter TodoFilter) []model.TodoItem
}

type TodoFilter struct {
	FromDate *time.Time //nil if empty
	ToDate   *time.Time
	Status   string
}

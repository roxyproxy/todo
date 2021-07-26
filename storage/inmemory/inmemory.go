package inmemory

import (
	"github.com/satori/go.uuid"
	"strings"
	"time"
	"todo/model"
	"todo/storage"
)

type InMemory struct {
	todoItems map[string]model.TodoItem
}

func NewInMemoryStorage() *InMemory {
	return &InMemory{map[string]model.TodoItem{}}
}

func (i *InMemory) GetItem(id string) (model.TodoItem, error) {
	todo := i.todoItems[id]
	return todo, nil
}

func (i *InMemory) UpdateItem(item model.TodoItem) error {
	i.todoItems[item.Id] = item
	return nil
}

func (i *InMemory) DeleteItem(id string) error {
	delete(i.todoItems, id)
	return nil
}

func (i *InMemory) AddItem(item model.TodoItem) (string, error) {
	u := uuid.NewV4().String()
	item.Id = u
	if item.Status == "" {
		item.Status = "new"
	}
	i.todoItems[u] = item
	return u, nil
}

func (i *InMemory) GetAllItems(filter storage.TodoFilter) ([]model.TodoItem, error) {
	arr := make([]model.TodoItem, 0)
	for _, value := range i.todoItems {
		if filter == (storage.TodoFilter{}) {
			arr = append(arr, value)
		} else {
			if filtered(filter, value) {
				arr = append(arr, value)
			}
		}
	}
	return arr, nil
}

func filtered(filter storage.TodoFilter, t model.TodoItem) bool {
	return statusOk(filter.Status, t.Status) && toDateOk(filter.ToDate, t.Date) && fromDateOk(filter.FromDate, t.Date)
}

func statusOk(status string, s string) bool {
	if status != "" && strings.Compare(status, s) != 0 {
		return false
	}
	return true
}

func toDateOk(toDate *time.Time, d time.Time) bool {
	if toDate != nil && toDate.Before(d) {
		return false
	}
	return true
}
func fromDateOk(fromDate *time.Time, d time.Time) bool {
	if fromDate != nil && fromDate.After(d) {
		return false
	}
	return true
}

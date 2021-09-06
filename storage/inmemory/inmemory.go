package inmemory

import (
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
	"todo/model"
	"todo/storage"
)

type InMemory struct {
	todoItems map[string]model.TodoItem
	users     map[string]model.User
}

func NewInMemoryStorage() *InMemory {
	return &InMemory{map[string]model.TodoItem{}, map[string]model.User{}}
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
			if itemFiltered(filter, value) {
				arr = append(arr, value)
			}
		}
	}
	return arr, nil
}

func itemFiltered(filter storage.TodoFilter, t model.TodoItem) bool {
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

// Users.
func (i *InMemory) GetUser(id string) (model.User, error) {
	user := i.users[id]
	return user, nil
}

func (i *InMemory) UpdateUser(u model.User) error {
	i.users[u.Id] = u
	return nil
}

func (i *InMemory) DeleteUser(id string) error {
	delete(i.users, id)
	return nil
}

func (i *InMemory) AddUser(user model.User) (string, error) {
	u := uuid.NewV4().String()
	user.Id = u
	i.users[u] = user
	return u, nil
}

func (i *InMemory) GetAllUsers(filter storage.UserFilter) ([]model.User, error) {
	arr := make([]model.User, 0)
	for _, value := range i.users {
		if filter == (storage.UserFilter{}) {
			arr = append(arr, value)
		} else {
			if userFiltered(filter, value) {
				arr = append(arr, value)
			}
		}
	}
	return arr, nil
}

func userFiltered(filter storage.UserFilter, t model.User) bool {
	return usernameOk(filter.UserName, t.UserName)
}

func usernameOk(username string, s string) bool {
	if username != "" && strings.Compare(username, s) != 0 {
		return false
	}
	return true
}

package inmemory

import (
	"fmt"
	"strings"
	"time"
	"todo/model"
	"todo/storage"
)

// InMemory represents in memory structure.
type InMemory struct {
	todoItems map[string]model.TodoItem
	users     map[string]model.User
}

// NewInMemoryStorage returns InMemory struct.
func NewInMemoryStorage() *InMemory {
	return &InMemory{map[string]model.TodoItem{}, map[string]model.User{}}
}

// GetItem gets item from memory.
func (i *InMemory) GetItem(id string) (model.TodoItem, error) {
	todo := i.todoItems[id]
	location, err := time.LoadLocation(i.users[todo.UserID].Location.String())
	if err != nil {
		return model.TodoItem{}, fmt.Errorf("cant load location")
	}
	todo.Date = todo.Date.In(location)
	return todo, nil
}

// UpdateItem updates todo in memory.
func (i *InMemory) UpdateItem(item model.TodoItem) error {
	if item.Date.IsZero() {
		item.Date = time.Now().UTC()
	} else {
		item.Date = item.Date.UTC()
	}
	i.todoItems[item.ID] = item
	return nil
}

// DeleteItem deletes todo from memory.
func (i *InMemory) DeleteItem(id string) error {
	delete(i.todoItems, id)
	return nil
}

// AddItem adds todo to memory.
func (i *InMemory) AddItem(item model.TodoItem) (string, error) {
	u := uuid.NewV4().String()
	item.ID = u
	if item.Status == "" {
		item.Status = "new"
	}

	if item.Date.IsZero() {
		item.Date = time.Now().UTC()
	} else {
		item.Date = item.Date.UTC()
	}

	i.todoItems[u] = item
	return u, nil
}

// GetAllItems gets all todos from memory.
func (i *InMemory) GetAllItems(filter storage.TodoFilter) ([]model.TodoItem, error) {
	arr := make([]model.TodoItem, 0)
	for _, value := range i.todoItems {
		location, err := time.LoadLocation(i.users[value.UserID].Location.String())
		if err != nil {
			return arr, fmt.Errorf("cant load location")
		}
		value.Date = value.Date.In(location)
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
	return useridOk(filter.UserID, t.UserID) && statusOk(filter.Status, t.Status) && toDateOk(filter.ToDate, t.Date) && fromDateOk(filter.FromDate, t.Date)
}

func useridOk(userid string, s string) bool {
	if userid != "" && strings.Compare(userid, s) != 0 {
		return false
	}
	return true
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

// GetUser gets user from memory.
func (i *InMemory) GetUser(id string) (model.User, error) {
	user := i.users[id]
	return user, nil
}

// UpdateUser updates user in memory.
func (i *InMemory) UpdateUser(u model.User) error {
	i.users[u.ID] = u
	return nil
}

// DeleteUser deletes user from memory.
func (i *InMemory) DeleteUser(id string) error {
	delete(i.users, id)
	return nil
}

// AddUser adds user to memory.
func (i *InMemory) AddUser(user model.User) (string, error) {
	u := uuid.NewV4().String()
	user.ID = u
	i.users[u] = user
	return u, nil
}

// GetAllUsers gets all users from memory.
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

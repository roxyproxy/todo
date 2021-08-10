package inmemory

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
	"todo/model"
	"todo/storage"
)

func TestStorage(t *testing.T) {
	l, _ := time.LoadLocation("America/New_York")
	location := model.CustomLocation{l}
	storageInMemory := InMemory{
		map[string]model.TodoItem{
			"6ba7b810-9dad-11d1-80b4-00c04fd430c8": model.TodoItem{
				Id:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
				Name: "todo1"},
		},
		map[string]model.User{
			"6ba7b810-9dad-11d1-80b4-00c04fd430c8": model.User{
				Id:        "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
				UserName:  "RoxyProxy",
				FirstName: "Roxy",
				LastName:  "Proxy",
				Password:  "$2a$14$Vv0FoIWcwWSf0mXMy.jFXebqBj/KXBetgN725ComfazcNemFUMVli",
				Location:  location,
			},
		},
	}

	t.Run("Get todo item", func(t *testing.T) {
		want := "todo1"
		todo, err := storageInMemory.GetItem("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		got := todo.Name

		if err != nil {
			t.Errorf("Error in GetItem %q", err)
		}

		assertEqual(t, got, want)
	})

	t.Run("Update todo item", func(t *testing.T) {
		storageInMemory.UpdateItem(model.TodoItem{Id: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", Name: "todo2"})

		want := "todo2"
		todo, err := storageInMemory.GetItem("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		got := todo.Name

		if err != nil {
			t.Errorf("Error in UpdateItem %q", err)
		}
		assertEqual(t, got, want)
	})

	t.Run("Delete todo item", func(t *testing.T) {
		storageInMemory.DeleteItem("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		want := model.TodoItem{}
		got, err := storageInMemory.GetItem("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

		if err != nil {
			t.Errorf("Error in DeleteItem %q", err)
		}
		assertEqual(t, got, want)
	})

	t.Run("Add todo item", func(t *testing.T) {
		todo := model.TodoItem{Name: "todo3"}
		got, err := storageInMemory.AddItem(todo)
		if err != nil {
			t.Errorf("Error in AddItem %q", err)
		}
		want, err := uuid.FromString(got)
		if err != nil {
			t.Errorf("got %+v, want %+v", got, want)
		}

		storageInMemory.DeleteItem(got)

	})

	t.Run("Get all todo items", func(t *testing.T) {
		todo1, _ := storageInMemory.AddItem(model.TodoItem{Name: "todo1"})
		todo2, _ := storageInMemory.AddItem(model.TodoItem{Name: "todo2"})

		todoitems, err := storageInMemory.GetAllItems(storage.TodoFilter{})
		if err != nil {
			t.Errorf("Error in GetAllItems %q", err)
		}

		got := len(todoitems)
		want := 2

		assertEqual(t, got, want)

		storageInMemory.DeleteItem(todo1)
		storageInMemory.DeleteItem(todo2)
	})

	t.Run("Get filtered todo items", func(t *testing.T) {
		date1 := time.Now()
		date2 := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
		date3 := time.Date(2029, 11, 17, 20, 34, 58, 651387237, time.UTC)
		todo1, _ := storageInMemory.AddItem(model.TodoItem{Name: "todo1", Date: date1, Status: "New"})
		todo2, _ := storageInMemory.AddItem(model.TodoItem{Name: "todo2", Date: date2, Status: "New"})
		todo3, _ := storageInMemory.AddItem(model.TodoItem{Name: "todo2", Date: date3, Status: "New"})

		date4 := time.Date(2000, 11, 17, 20, 34, 58, 651387237, time.UTC)

		filter := storage.TodoFilter{ToDate: &date1, Status: "New", FromDate: &date4}
		todoitems, err := storageInMemory.GetAllItems(filter)

		if err != nil {
			t.Errorf("Error in GetAllItems %q", err)
		}

		got := len(todoitems)
		want := 2

		assertEqual(t, got, want)
		storageInMemory.DeleteItem(todo1)
		storageInMemory.DeleteItem(todo2)
		storageInMemory.DeleteItem(todo3)

	})

	t.Run("Get user", func(t *testing.T) {
		l, _ := time.LoadLocation("America/New_York")
		location := model.CustomLocation{l}
		want := model.User{Id: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			UserName:  "RoxyProxy",
			FirstName: "Roxy",
			LastName:  "Proxy",
			Password:  "$2a$14$Vv0FoIWcwWSf0mXMy.jFXebqBj/KXBetgN725ComfazcNemFUMVli",
			Location:  location}
		user, err := storageInMemory.GetUser("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

		if err != nil {
			t.Errorf("Error in GetUser %q", err)
		}
		assert.Equal(t, want, user)

		storageInMemory.DeleteUser("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	})
	t.Run("Add user", func(t *testing.T) {
		newUser := model.User{UserName: "Roxy1", Password: "Proxy1"}
		id, err := storageInMemory.AddUser(newUser)
		if err != nil {
			t.Errorf("Error in AddUser %q", err)
		}

		user, _ := storageInMemory.GetUser(id)
		newUser.Id = id
		assert.Equal(t, newUser, user)

		storageInMemory.DeleteUser(id)
	})

	t.Run("Delete user", func(t *testing.T) {
		newUser := model.User{UserName: "Roxy2", Password: "Proxy2"}
		id, _ := storageInMemory.AddUser(newUser)
		err := storageInMemory.DeleteUser(id)
		if err != nil {
			t.Errorf("Error in DeleteUser %q", err)
		}
		user, _ := storageInMemory.GetUser(id)
		assert.Equal(t, user, model.User{})

	})

	t.Run("Update user", func(t *testing.T) {
		newUser := model.User{UserName: "Roxy2", Password: "Proxy2"}
		id, _ := storageInMemory.AddUser(newUser)
		newUser.Id = id
		newUser.UserName = "Roxy3"
		err := storageInMemory.UpdateUser(newUser)
		if err != nil {
			t.Errorf("Error in UpdateUser %q", err)
		}
		user, _ := storageInMemory.GetUser(id)
		assert.Equal(t, newUser.UserName, user.UserName)

		storageInMemory.DeleteUser(id)
	})

	t.Run("Get all users", func(t *testing.T) {
		user1 := model.User{UserName: "Roxy1", Password: "Proxy1"}
		user2 := model.User{UserName: "Roxy2", Password: "Proxy2"}
		user3 := model.User{UserName: "Roxy3", Password: "Proxy3"}
		id1, _ := storageInMemory.AddUser(user1)
		id2, _ := storageInMemory.AddUser(user2)
		id3, _ := storageInMemory.AddUser(user3)

		users, err := storageInMemory.GetAllUsers(storage.UserFilter{})
		if err != nil {
			t.Errorf("Error in GetAllUsers %q", err)
		}

		assert.Equal(t, 3, len(users))

		storageInMemory.DeleteUser(id1)
		storageInMemory.DeleteUser(id2)
		storageInMemory.DeleteUser(id3)

	})

	t.Run("Get filtered users", func(t *testing.T) {
		user1 := model.User{UserName: "Roxy1", Password: "Proxy1"}
		user2 := model.User{UserName: "Roxy2", Password: "Proxy2"}
		user3 := model.User{UserName: "Roxy1", Password: "Proxy1"}
		id1, _ := storageInMemory.AddUser(user1)
		id2, _ := storageInMemory.AddUser(user2)
		id3, _ := storageInMemory.AddUser(user3)

		users, err := storageInMemory.GetAllUsers(storage.UserFilter{UserName: "Roxy1"})
		if err != nil {
			t.Errorf("Error in GetAllUsers %q", err)
		}

		assert.Equal(t, 2, len(users))

		storageInMemory.DeleteUser(id1)
		storageInMemory.DeleteUser(id2)
		storageInMemory.DeleteUser(id3)
	})

}

func assertEqual(t *testing.T, got interface{}, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

package inmemory

import (
	uuid "github.com/satori/go.uuid"
	"reflect"
	"testing"
	"time"
	"todo/model"
	"todo/storage"
)

func TestStorage(t *testing.T) {
	storageInMemory := InMemory{
		map[string]model.TodoItem{
			"6ba7b810-9dad-11d1-80b4-00c04fd430c8": model.TodoItem{
				Id:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
				Name: "todo1"},
		},
	}

	t.Run("Get todo item", func(t *testing.T) {
		want := "todo1"
		todo, err := storageInMemory.GetItem("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		got := todo.Name

		if err != nil {
			t.Errorf("Error in GetItem %q", err)
		}

		assertTodoItem(t, got, want)
	})

	t.Run("Update todo item", func(t *testing.T) {
		storageInMemory.UpdateItem(model.TodoItem{Id: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", Name: "todo2"})

		want := "todo2"
		todo, err := storageInMemory.GetItem("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		got := todo.Name

		if err != nil {
			t.Errorf("Error in UpdateItem %q", err)
		}
		assertTodoItem(t, got, want)
	})

	t.Run("Delete todo item", func(t *testing.T) {
		storageInMemory.DeleteItem("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		want := model.TodoItem{}
		got, err := storageInMemory.GetItem("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

		if err != nil {
			t.Errorf("Error in DeleteItem %q", err)
		}
		assertTodoItem(t, got, want)
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

		assertTodoItem(t, got, want)

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

		filter := storage.TodoFilter{ToDate: &date1, Status: "New"}
		todoitems, err := storageInMemory.GetAllItems(filter)

		if err != nil {
			t.Errorf("Error in GetAllItems %q", err)
		}

		got := len(todoitems)
		want := 2

		assertTodoItem(t, got, want)
		storageInMemory.DeleteItem(todo1)
		storageInMemory.DeleteItem(todo2)
		storageInMemory.DeleteItem(todo3)

	})

}

func assertTodoItem(t *testing.T, got interface{}, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

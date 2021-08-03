package server

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"todo/model"
	mockstore "todo/server/mocks"
	"todo/storage"
	"todo/storage/inmemory"
)

func TestServerWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockstore.NewMockStorage(ctrl)
	server := NewTodoServer(m)

	t.Run("get all items (empty)", func(t *testing.T) {
		m.EXPECT().GetAllItems(storage.TodoFilter{}).Return([]model.TodoItem{}, nil)

		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/todos", nil)
		server.Serve.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.JSONEq(t, "[]", response.Body.String())
	})

	t.Run("get filtered items", func(t *testing.T) {
		todoItems := []model.TodoItem{
			{"123", "test1", time.Now(), "done"},
			{"124", "test2", time.Now(), "done"},
		}
		todoItemsJson, _ := json.Marshal(&todoItems)
		filter := storage.TodoFilter{Status: "done"}
		m.EXPECT().GetAllItems(filter).Return(todoItems, nil)

		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/todos?status=done", nil)
		server.Serve.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json")
		assert.JSONEq(t, string(todoItemsJson), response.Body.String())
	})

	t.Run("add new item", func(t *testing.T) {
		todo := model.TodoItem{Name: "test1"}
		todoId := model.TodoId{Id: "123"}
		m.EXPECT().AddItem(todo).Return("123", nil)

		todoJson, _ := json.Marshal(&todo)
		todoIdJson, _ := json.Marshal(&todoId)

		request, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(todoJson))
		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.JSONEq(t, string(todoIdJson), response.Body.String())
	})

}

func TestServer(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()
	server := NewTodoServer(storage)
	assert := assert.New(t)

	t.Run("add new item", func(t *testing.T) {
		todo := model.TodoItem{Name: "test1"}
		j, _ := json.Marshal(&todo)

		resp := AddRequest(server, j)
		assertStatus(t, resp.Code, http.StatusOK)

		todoId := model.TodoId{}
		json.NewDecoder(resp.Body).Decode(&todoId)

		want, err := uuid.FromString(todoId.Id)
		if err != nil {
			t.Errorf("got %+v, want %+v", todoId.Id, want)
		}

		DeleteTodo(server, todoId.Id)

	})

	t.Run("get not existing item", func(t *testing.T) {
		v := url.Values{}
		v.Set("id", "test")
		assert.HTTPStatusCode(server.getItemHandler, "GET", "/todo", v, 404)
	})

	t.Run("delete item", func(t *testing.T) {
		todo := model.TodoItem{Name: "test1"}
		j, _ := json.Marshal(&todo)

		resp := AddRequest(server, j)

		todoId := model.TodoId{}
		json.NewDecoder(resp.Body).Decode(&todoId)

		req, _ := http.NewRequest(http.MethodDelete, "/todos/"+todoId.Id, nil)
		res := httptest.NewRecorder()
		server.Serve.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
	})

	t.Run("update item", func(t *testing.T) {
		todo1 := model.TodoItem{Name: "test1"}
		j1, _ := json.Marshal(&todo1)
		todo2 := model.TodoItem{Name: "test2"}
		j2, _ := json.Marshal(&todo2)
		todo3 := model.TodoItem{}

		resp := AddRequest(server, j1)
		todoId := model.TodoId{}
		json.NewDecoder(resp.Body).Decode(&todoId)

		req, _ := http.NewRequest(http.MethodPut, "/todos/"+todoId.Id, bytes.NewBuffer(j2))
		resp = httptest.NewRecorder()
		server.Serve.ServeHTTP(resp, req)
		json.NewDecoder(resp.Body).Decode(&todo3)

		assertStatus(t, resp.Code, http.StatusOK)

		DeleteTodo(server, todoId.Id)

	})

	t.Run("get all items", func(t *testing.T) {
		todo := model.TodoItem{Name: "test12"}
		j, _ := json.Marshal(&todo)
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(j))
		server.Serve.ServeHTTP(resp, req)

		var err error
		req, err = http.NewRequest(http.MethodGet, "/todos", nil)
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		server.Serve.ServeHTTP(res, req)
		todoItems := []model.TodoItem{}

		json.NewDecoder(res.Body).Decode(&todoItems)

		assertStatus(t, res.Code, http.StatusOK)
		assertResponseBody(t, len(todoItems), 1)
	})

}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
func assertResponseBody(t testing.TB, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func AddRequest(server *TodoServer, b []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(b))
	resp := httptest.NewRecorder()
	server.Serve.ServeHTTP(resp, req)

	return resp
}

func DeleteTodo(server *TodoServer, id string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodDelete, "/todos/"+id, nil)
	res := httptest.NewRecorder()
	server.Serve.ServeHTTP(res, req)

	return res
}

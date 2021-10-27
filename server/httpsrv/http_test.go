package httpsrv

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	conf "todo/config"
	"todo/logger"
	"todo/model"
	mockstore "todo/server/mocks"
	"todo/service"
	"todo/storage"
)

func TestServerWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := conf.New()
	m := mockstore.NewMockStorage(ctrl)
	s := service.NewService(m, c)

	server := NewHTTPServer(s, c, logger.New(ioutil.Discard))

	username := "Roxy"
	password := "SecretPassword12!"
	l, _ := time.LoadLocation("America/New_York")
	location := model.CustomLocation{Location: l}
	credentials := model.Credentials{UserName: username, Password: password}
	hash, err := server.service.HashPassword(password)
	assert.NoError(t, err)
	user := model.User{
		ID:        "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		UserName:  username,
		FirstName: "Roxy",
		LastName:  "Proxy",
		Password:  hash,
		Location:  location,
	}
	token, _ := server.service.GenerateToken(user.ID, server.config.SecretKey)

	t.Run("test hashPassword", func(t *testing.T) {
		_, err := server.service.HashPassword(password)
		assert.NoError(t, err)
	})

	t.Run("test checkPassword", func(t *testing.T) {
		hash, err := server.service.HashPassword(password)
		assert.NoError(t, err)

		ok := server.service.CheckPasswordHash(password, hash)
		assert.True(t, ok)
	})

	t.Run("test authenticate user", func(t *testing.T) {
		users := []model.User{user}
		m.EXPECT().GetAllUsers(storage.UserFilter{UserName: credentials.UserName}).Return(users, nil)

		_, err := server.service.AuthenticateUser(credentials)
		assert.NoError(t, err)
	})

	t.Run("test generate Token", func(t *testing.T) {
		_, err := server.service.GenerateToken(user.ID, server.config.SecretKey)
		assert.NoError(t, err)
	})

	t.Run("test get Token", func(t *testing.T) {
		tkn, _ := server.service.GenerateToken(user.ID, server.config.SecretKey)
		r, _ := http.NewRequest(http.MethodGet, "/users", nil)
		r.Header.Set("Authorization", "Bearer "+tkn.TokenString)

		got := getToken(r)
		assert.NotEmpty(t, got)
	})

	t.Run("test login user", func(t *testing.T) {
		users := []model.User{user}
		m.EXPECT().GetAllUsers(storage.UserFilter{UserName: credentials.UserName}).Return(users, nil)

		credentialsJSON, err := json.Marshal(&credentials)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(credentialsJSON))
		assert.NoError(t, err)

		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	})

	//todo check all fields except password
	/*
		t.Run("add user", func(t *testing.T) {
			newUser := model.User{
				UserName:  user.UserName,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Location:  user.Location,
				Password:  password,
			}
			userID := model.TodoID{ID: user.ID}

			m.EXPECT().
				AddUser(newUser).
				DoAndReturn(func(u model.User) (string, error) {
					if !checkPasswordHash(u.Password, user.Password) {
						t.Fail()
					}
					u.Password = "test"
					fmt.Println(u)
					return user.ID, nil
				})

			userJson, err := json.Marshal(&newUser)
			assert.NoError(t, err)

			userIDJson, err := json.Marshal(&userID)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJson))
			assert.NoError(t, err)

			response := httptest.NewRecorder()
			server.Serve.ServeHTTP(response, request)
			assert.Equal(t, response.Code, http.StatusOK)
			assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
			assert.JSONEq(t, string(userIDJson), response.Body.String())

		})
	*/

	t.Run("get user", func(t *testing.T) {
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().GetUser(user.ID).Return(user, nil)

		userJSON, err := json.Marshal(&user)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodGet, "/users/"+user.ID, bytes.NewBuffer(userJSON))
		assert.NoError(t, err)

		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json")
		assert.JSONEq(t, string(userJSON), response.Body.String())
	})

	t.Run("get all users", func(t *testing.T) {
		users := []model.User{user}
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().GetAllUsers(storage.UserFilter{}).Return(users, nil)

		userJSON, err := json.Marshal(&users)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodGet, "/users", bytes.NewBuffer(userJSON))
		assert.NoError(t, err)

		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json")
		assert.JSONEq(t, string(userJSON), response.Body.String())
	})

	t.Run("get all users filtered", func(t *testing.T) {
		users := []model.User{user}
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().GetAllUsers(storage.UserFilter{UserName: user.UserName}).Return(users, nil)

		userJSON, err := json.Marshal(&users)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodGet, "/users?username=Roxy", bytes.NewBuffer(userJSON))
		assert.NoError(t, err)

		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json")
		assert.JSONEq(t, string(userJSON), response.Body.String())
	})

	t.Run("update user", func(t *testing.T) {
		newuser := model.User{
			ID:        user.ID,
			UserName:  "Roxy2",
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Location:  user.Location,
		}
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().UpdateUser(newuser).Return(nil)

		userJSON, err := json.Marshal(&newuser)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPut, "/users/"+user.ID, bytes.NewBuffer(userJSON))
		assert.NoError(t, err)

		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json")
	})

	t.Run("delete user", func(t *testing.T) {
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().DeleteUser(user.ID).Return(nil)

		request, err := http.NewRequest(http.MethodDelete, "/users/"+user.ID, nil)
		assert.NoError(t, err)

		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json")
	})

	t.Run("get not existing item", func(t *testing.T) {
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().GetItem("123").Return(model.TodoItem{}, nil)

		request, err := http.NewRequest(http.MethodGet, "/todos/123", nil)
		assert.NoError(t, err)

		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, http.StatusForbidden, response.Code)
	})

	t.Run("get all items unathorized", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/todos", nil)
		assert.NoError(t, err)

		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("get all items (empty)", func(t *testing.T) {
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		m.EXPECT().GetAllItems(storage.TodoFilter{UserID: user.ID}).Return([]model.TodoItem{}, nil)

		response := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/todos", nil)
		assert.NoError(t, err)

		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.JSONEq(t, "[]", response.Body.String())
	})

	t.Run("get filtered items", func(t *testing.T) {
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		todoItems := []model.TodoItem{
			{ID: "123", Name: "test1", Date: time.Now(), Status: "done"},
			{ID: "124", Name: "test2", Date: time.Now(), Status: "done"},
		}
		todoItemsJSON, err := json.Marshal(&todoItems)
		assert.NoError(t, err)

		filter := storage.TodoFilter{UserID: user.ID, Status: "done"}

		m.EXPECT().GetAllItems(filter).Return(todoItems, nil)

		response := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/todos?status=done", nil)
		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		assert.NoError(t, err)

		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, response.Header().Get("Content-Type"), "application/json")
		assert.JSONEq(t, string(todoItemsJSON), response.Body.String())
	})

	t.Run("add new item", func(t *testing.T) {
		m.EXPECT().GetUser(user.ID).Return(user, nil)
		todo := model.TodoItem{Name: "test1", UserID: user.ID}
		todoID := model.TodoID{ID: "123"}
		m.EXPECT().AddItem(todo).Return("123", nil)

		todoJSON, err := json.Marshal(&todo)
		assert.NoError(t, err)

		todoIDJSON, err := json.Marshal(&todoID)
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(todoJSON))
		request.Header.Set("Authorization", "Bearer "+token.TokenString)
		assert.NoError(t, err)

		response := httptest.NewRecorder()
		server.Serve.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
		assert.JSONEq(t, string(todoIDJSON), response.Body.String())
	})
}

/*
func TestServer(t *testing.T) {
	storage := inmemory.NewInMemoryStorage()
	server := NewTodoServer(storage, conf.New(), logger.New(ioutil.Discard))

	username := "Roxy"
	password := "SecretPassword12!"
	l, _ := time.LoadLocation("America/New_York")
	location := model.CustomLocation{Location: l}
	hash, err := hashPassword(password)
	if err != nil {
		t.Errorf("hashPassword %+v", err)
	}
	user := model.User{
		ID:        "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		UserName:  username,
		FirstName: "Roxy",
		LastName:  "Proxy",
		Password:  hash,
		Location:  location,
	}
	token, _ := generateToken(user.ID, server.config.SecretKey)

	t.Run("add new item", func(t *testing.T) {
		todo := model.TodoItem{Name: "test1"}
		j, _ := json.Marshal(&todo)

		resp := AddRequest(server, j)
		assertStatus(t, resp.Code, http.StatusOK)

		todoID := model.TodoID{}
		err := json.NewDecoder(resp.Body).Decode(&todoID)
		assert.NoError(t, err)

		want, err := uuID.FromString(todoID.ID)
		if err != nil {
			t.Errorf("got %+v, want %+v", todoID.ID, want)
		}

		DeleteTodo(server, todoID.ID)
	})

	t.Run("delete item", func(t *testing.T) {
		todo := model.TodoItem{Name: "test1"}
		j, _ := json.Marshal(&todo)

		resp := AddRequest(server, j)
		todoID := model.TodoID{}
		err := json.NewDecoder(resp.Body).Decode(&todoID)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodDelete, "/todos/"+todoID.ID, nil)
		req.Header.Set("Authorization", "Bearer "+token.TokenString)

		res := httptest.NewRecorder()
		server.Serve.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusOK)
	})

	t.Run("update item", func(t *testing.T) {
		todo1 := model.TodoItem{Name: "test1"}
		j1, _ := json.Marshal(&todo1)
		todo2 := model.TodoItem{Name: "test2"}
		j2, _ := json.Marshal(&todo2)

		resp := AddRequest(server, j1)
		todoID := model.TodoID{}
		err := json.NewDecoder(resp.Body).Decode(&todoID)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPut, "/todos/"+todoID.ID, bytes.NewBuffer(j2))
		req.Header.Set("Authorization", "Bearer "+token.TokenString)

		resp = httptest.NewRecorder()
		server.Serve.ServeHTTP(resp, req)
		assertStatus(t, resp.Code, http.StatusOK)

		DeleteTodo(server, todoID.ID)
	})

	t.Run("get all items", func(t *testing.T) {
		todo := model.TodoItem{Name: "test12"}
		j, _ := json.Marshal(&todo)
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(j))
		req.Header.Set("Authorization", "Bearer "+token.TokenString)

		server.Serve.ServeHTTP(resp, req)

		var err error
		req, err = http.NewRequest(http.MethodGet, "/todos", nil)
		req.Header.Set("Authorization", "Bearer "+token.TokenString)

		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		server.Serve.ServeHTTP(res, req)
		todoItems := []model.TodoItem{}

		err = json.NewDecoder(res.Body).Decode(&todoItems)
		assert.NoError(t, err)

		assertStatus(t, res.Code, http.StatusOK)
		assertResponseBody(t, len(todoItems), 1)
	})
}



func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("dID not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func AddRequest(server *TodoServer, b []byte) *httptest.ResponseRecorder {
	token, _ := generateToken("6ba7b810-9dad-11d1-80b4-00c04fd430c8", server.config.SecretKey)
	req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+token.TokenString)
	resp := httptest.NewRecorder()
	server.Serve.ServeHTTP(resp, req)

	return resp
}

func DeleteTodo(server *TodoServer, ID string) *httptest.ResponseRecorder {
	token, _ := generateToken("6ba7b810-9dad-11d1-80b4-00c04fd430c8", server.config.SecretKey)
	req, _ := http.NewRequest(http.MethodDelete, "/todos/"+ID, nil)
	req.Header.Set("Authorization", "Bearer "+token.TokenString)

	res := httptest.NewRecorder()
	server.Serve.ServeHTTP(res, req)

	return res
}


*/

package integration_test

import (
	"bytes"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"os"
	"todo/logger"

	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"

	"context"
	"encoding/json"
	"testing"
	"time"
	conf "todo/config"
	"todo/model"
)

func TestUsers(t *testing.T) {
	log := logger.New(os.Stderr)
	if err := godotenv.Load("../../.env"); err != nil {
		log.Warning("No .env file found")
	}

	config := conf.New()

	dbpool, err := pgxpool.Connect(context.Background(), config.DBUrl)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	l, _ := time.LoadLocation("America/New_York")
	newUser := model.NewUser{
		UserName:  "RoxyProxy",
		FirstName: "Roxy",
		LastName:  "Proxy",
		Password:  "RoxyP",
		Location:  model.CustomLocation{Location: l},
	}

	type User struct {
		ID        string `db:"id"`
		UserName  string `db:"username"`
		FirstName string `db:"firstname"`
		LastName  string `db:"lastname"`
		Password  string `db:"password"`
		Location  string `db:"location"`
	}
	var userStruct = sqlbuilder.NewStruct(new(User)).For(sqlbuilder.PostgreSQL)
	userID := uuid.NewV4().String()
	user := User{
		ID:        userID,
		UserName:  "RoxyProxy2",
		FirstName: "Roxy",
		LastName:  "Proxy",
		Password:  "RoxyP",
		Location:  "America/New_York",
	}

	ib := userStruct.InsertInto("users", &user)

	// Execute the query.
	sql, args := ib.Build()

	_, err = dbpool.Exec(context.Background(), sql, args...)
	assert.NoError(t, err)

	t.Run("add user", func(t *testing.T) {
		userJson, err := json.Marshal(&newUser)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost%v/users", config.HTTPPort), bytes.NewBuffer(userJson))

		assert.NoError(t, err)

		client := http.Client{}
		response, err := client.Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
		userID := model.TodoID{}
		err = json.NewDecoder(response.Body).Decode(&userID)
		assert.NoError(t, err)
		ID, err := uuid.FromString(userID.ID)
		assert.IsType(t, uuid.UUID{}, ID)
	})

	t.Run("user login", func(t *testing.T) {
		userJson, err := json.Marshal(&newUser)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost%v/user/login", config.HTTPPort), bytes.NewBuffer(userJson))

		assert.NoError(t, err)

		client := http.Client{}
		response, err := client.Do(req)

		assert.NoError(t, err)
		defer response.Body.Close()

		fmt.Println(response.Body)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
		token := model.Token{}
		err = json.NewDecoder(response.Body).Decode(&token)
		assert.NoError(t, err)

	})

}

package db

import (
	"context"
	"fmt"
	"reflect"
	"time"
	"todo/storage"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
	"todo/model"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(p *pgxpool.Pool) *Postgres {
	return &Postgres{pool: p}
}

func (i *Postgres) AddUser(user model.User) (string, error) {
	u := uuid.NewV4().String()
	user.Id = u

	err := i.pool.QueryRow(context.Background(),
		"INSERT INTO users (id, username, firstname, lastname, password, location) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Id, user.UserName, user.FirstName, user.LastName, user.Password, user.Location.String()).Scan(&user.Id)

	if err != nil {
		return "", fmt.Errorf("Unable to INSERT: %v", err)
	}
	return u, nil
}

func (i *Postgres) GetUser(id string) (model.User, error) {
	var user = model.User{}
	var l string
	err := i.pool.QueryRow(context.Background(),
		"SELECT id, username, firstname, lastname, password, location FROM users WHERE id = $1",
		id).Scan(&user.Id, &user.UserName, &user.FirstName, &user.LastName, &user.Password, &l)
	if err == pgx.ErrNoRows {
		return model.User{}, nil
	}

	if err != nil {
		return model.User{}, fmt.Errorf("Unable to SELECT: %v", err)
	}

	location, err := time.LoadLocation(l)
	if err != nil {
		return model.User{}, fmt.Errorf("Unable to convert location: %v", err)
	}
	user.Location = model.CustomLocation{location}

	if err != nil {
		return model.User{}, fmt.Errorf("Unable to SELECT: %v", err)
	}
	return user, nil
}

func (i *Postgres) UpdateUser(u model.User) error {
	_, err := i.pool.Exec(context.Background(),
		"UPDATE users SET username = $2, firstname = $3, lastname=$4, password=$5, location=$6 WHERE id = $1",
		u.Id, u.UserName, u.FirstName, u.LastName, u.Password, u.Location.String())

	if err != nil {
		return fmt.Errorf("Unable to UPDATE: %v\n", err)
	}

	/*if ct.RowsAffected() == 0 {
		return fmt.Errorf("User not found: %v\n", err)
	}*/

	return nil
}

func (i *Postgres) DeleteUser(id string) error {
	_, err := i.pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("Unable to DELETE: %v", err)
	}

	/*if ct.RowsAffected() == 0 {
		return fmt.Errorf("User not found: %v\n", err)
	}*/
	return nil
}

func (i *Postgres) GetAllUsers(filter storage.UserFilter) ([]model.User, error) {
	arr := make([]model.User, 0)
	query := "SELECT id, username, firstname, lastname, password, location FROM users WHERE 1=1"
	v := reflect.ValueOf(filter)

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i).Interface().(string)
		if len(val) > 0 {
			query += " and " + v.Type().Field(i).Name + "=" + "'" + val + "'"
		}
	}

	rows, err := i.pool.Query(context.Background(), query)
	if err == pgx.ErrNoRows {
		return arr, nil
	}
	if err != nil {
		return arr, fmt.Errorf("Unable to SELECT: %v", err)
	}

	var l string
	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.Id, &user.UserName, &user.FirstName, &user.LastName, &user.Password, &l)

		if err != nil {
			return arr, fmt.Errorf("Unable to SELECT: %v", err)
		}
		location, err := time.LoadLocation(l)
		if err != nil {
			return arr, fmt.Errorf("Unable to convert location: %v", err)
		}
		user.Location = model.CustomLocation{location}

		arr = append(arr, user)
	}

	return arr, nil
}

func (i *Postgres) GetItem(id string) (model.TodoItem, error) {
	var todo = model.TodoItem{}

	err := i.pool.QueryRow(context.Background(),
		"SELECT id, name, date, status, userid FROM todos WHERE id = $1",
		id).Scan(&todo.Id, &todo.Name, &todo.Date, &todo.Status, &todo.UserId)

	if err == pgx.ErrNoRows {
		return model.TodoItem{}, nil
	}

	if err != nil {
		return model.TodoItem{}, fmt.Errorf("Unable to SELECT: %v", err)
	}

	var user = model.User{}
	var l string
	err = i.pool.QueryRow(context.Background(),
		"SELECT id, username, firstname, lastname, password, location FROM users WHERE id = $1",
		todo.UserId).Scan(&user.Id, &user.UserName, &user.FirstName, &user.LastName, &user.Password, &l)

	if err == pgx.ErrNoRows {
		return model.TodoItem{}, nil
	}
	if err != nil {
		return model.TodoItem{}, fmt.Errorf("Unable to SELECT: %v", err)
	}

	location, err := time.LoadLocation(l)
	if err != nil {
		return model.TodoItem{}, fmt.Errorf("cant load location")
	}
	todo.Date = todo.Date.In(location)

	return todo, nil
}

func (i *Postgres) UpdateItem(item model.TodoItem) error {
	if item.Date.IsZero() {
		item.Date = time.Now().UTC()
	} else {
		item.Date = item.Date.UTC()
	}

	_, err := i.pool.Exec(context.Background(),
		"UPDATE todos SET name=$2, date=$3, status=$4, userid=$5 WHERE id = $1",
		item.Id, item.Name, item.Date, item.Status, item.UserId)

	if err != nil {
		return fmt.Errorf("Unable to UPDATE: %v\n", err)
	}
	return nil
}

func (i *Postgres) DeleteItem(id string) error {
	_, err := i.pool.Exec(context.Background(), "DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("Unable to DELETE: %v", err)
	}
	return nil
}

func (i *Postgres) AddItem(item model.TodoItem) (string, error) {
	u := uuid.NewV4().String()
	item.Id = u
	if item.Status == "" {
		item.Status = "new"
	}

	if item.Date.IsZero() {
		item.Date = time.Now().UTC()
	} else {
		item.Date = item.Date.UTC()
	}

	err := i.pool.QueryRow(context.Background(),
		"INSERT INTO todos (id, name, date, status, userid) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		item.Id, item.Name, item.Date, item.Status, item.UserId).Scan(&item.Id)

	if err != nil {
		return "", fmt.Errorf("Unable to INSERT: %v", err)
	}
	return u, nil
}

func (i *Postgres) GetAllItems(filter storage.TodoFilter) ([]model.TodoItem, error) {
	arr := make([]model.TodoItem, 0)

	var user = model.User{}
	var l string
	err := i.pool.QueryRow(context.Background(),
		"SELECT id, username, firstname, lastname, password, location FROM users WHERE id = $1",
		filter.UserId).Scan(&user.Id, &user.UserName, &user.FirstName, &user.LastName, &user.Password, &l)
	if err != nil {
		return nil, fmt.Errorf("Unable to SELECT: %v", err)
	}

	location, err := time.LoadLocation(l)
	if err != nil {
		return arr, fmt.Errorf("cant load location")
	}

	query := "SELECT id, name, date, status, userid FROM todos WHERE 1=1"

	if len(filter.UserId) > 0 {
		query += " and userid = '" + filter.UserId + "'"
	}
	if len(filter.Status) > 0 {
		query += " and status = '" + filter.Status + "'"
	}
	if filter.FromDate != nil {
		query += " and date >= '" + filter.FromDate.UTC().Format(time.RFC3339) + "'"
	}
	if filter.ToDate != nil {
		query += " and date <= '" + filter.ToDate.UTC().Format(time.RFC3339) + "'"
	}

	rows, err := i.pool.Query(context.Background(), query)
	if err == pgx.ErrNoRows {
		return arr, nil
	}
	if err != nil {
		return arr, fmt.Errorf("Unable to SELECT: %v", err)
	}

	for rows.Next() {
		item := model.TodoItem{}
		err := rows.Scan(&item.Id, &item.Name, &item.Date, &item.Status, &item.UserId)
		if err != nil {
			return arr, fmt.Errorf("Unable to SELECT: %v", err)
		}
		item.Date = item.Date.In(location)
		arr = append(arr, item)
	}

	return arr, nil
}

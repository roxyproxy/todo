package model

import "time"

// TodoItem represents todo.
type TodoItem struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Status string    `json:"status"`
	UserID string    `json:"-"`
}

// TodoID represents todos id.
type TodoID struct {
	ID string `json:"id"`
}

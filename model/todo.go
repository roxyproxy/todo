package model

import "time"

type TodoItem struct {
	Id     string
	Name   string
	Date   time.Time //time.RFC3339
	Status string
}

type TodoId struct {
	Id string
}

package model

import "time"

type TodoItem struct {
	Id     string    `json:"id"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Status string    `json:"status"`
	UserId string    `json:"-"`
}

type TodoId struct {
	Id string `json:"id"`
}

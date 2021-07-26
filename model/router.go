package model

import (
	"net/http"
)

type Route struct {
	Method string
	Path   string
	Handlr http.HandlerFunc
}

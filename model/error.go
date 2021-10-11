package model

import (
	"errors"
)

// ErrOperational errors model defined for application.
var (
	ErrOperational  = errors.New("operational")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("forbidden")
)

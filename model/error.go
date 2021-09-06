package model

import (
	"errors"
)

var (
	ErrOperational  = errors.New("Operational")
	ErrBadRequest   = errors.New("Bad Request")
	ErrUnauthorized = errors.New("Unauthorized")
	ErrNotFound     = errors.New("Forbidden")
)

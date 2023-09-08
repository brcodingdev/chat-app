package errors

import (
	"errors"
)

var (
	ErrInCredentials = errors.New("invalid credentials")
	ErrInRequest     = errors.New("bad request")
	ErrDupEmail      = errors.New("email already exists")
	ErrToken         = errors.New("error jwt token")
)

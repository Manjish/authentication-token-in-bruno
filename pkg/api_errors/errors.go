package api_errors

import (
	"errors"
)

var (
	ErrUnauthorizedAccess   = errors.New("unauthorized access")
	ErrInvalidUUID          = errors.New("invalid uuid")
	ErrIncorrectCredentials = errors.New("incorrect credentials")
)

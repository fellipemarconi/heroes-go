package apierror

import "errors"

var (
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidBody        = errors.New("invalid request body")
)

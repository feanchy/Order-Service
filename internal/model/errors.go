package model

import "errors"

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrEmptyStr       = errors.New("status is required")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInternalServer = errors.New("internal server error")
)

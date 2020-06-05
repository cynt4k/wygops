package repository

import "errors"

var (
	// ErrAlreadyExists : Error if the entry does already exist
	ErrAlreadyExists = errors.New("already exists")
)

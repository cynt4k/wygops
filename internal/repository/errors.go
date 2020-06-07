package repository

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	// ErrAlreadyExists : Error if the entry does already exist
	ErrAlreadyExists = errors.New("already exists")
	// ErrNilID : Error when id is nil
	ErrNilID = errors.New("nil id")
	// ErrNotFound : Error when entry not found
	ErrNotFound = errors.New("not found")
	// ErrForbidden : Error when access is forbidden
	ErrForbidden = errors.New("forbidden")
)

func convertError(err error) error {
	switch {
	case gorm.IsRecordNotFoundError(err):
		return ErrNotFound
	default:
		return err
	}
}

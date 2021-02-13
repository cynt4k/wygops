package repository

import (
	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/pkg/util/optional"
)

// UserRepository : User repository to predifine the interfaces
type UserRepository interface {
	CreateUser(*models.User) (*models.User, error)
	GetUser(uint) (*models.User, error)
	GetUsers() (*[]models.User, error)
	GetLdapUsers() (*[]models.User, error)
	GetUserByUsername(string) (*models.User, error)
	GetUserWithDevices(uint) (*models.User, error)
	UpdateUser(userID uint, args UpdateUserArgs) (*models.User, error)
	DeleteUser(userID uint) error
}

type UpdateUserArgs struct {
	ProtectPassword optional.String
}

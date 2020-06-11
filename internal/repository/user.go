package repository

import "github.com/cynt4k/wygops/internal/models"

// UserRepository : User repository to predifine the interfaces
type UserRepository interface {
	CreateUser(*models.User) (*models.User, error)
	GetUser(uint) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
	GetUserWithDevices(uint) (*models.User, error)
	DeleteUser(userID uint) error
}

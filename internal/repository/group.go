package repository

import "github.com/cynt4k/wygops/internal/models"

// GroupRepository : Group repository to predefine the interface
type GroupRepository interface {
	CreateGroup(*models.Group) (*models.Group, error)
	AddUserToGroup(userID uint, groupID uint) error
	RemoveUserFromGroup(userID uint, groupID uint) error
	GetGroup(uint) (*models.Group, error)
	DeleteGroup(groupID uint) error
}

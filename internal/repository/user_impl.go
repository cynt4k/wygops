package repository

import (
	"github.com/cynt4k/wygops/internal/event"
	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/internal/util/gormutil"
	"github.com/jinzhu/gorm"
)

// CreateUser : Create an user
func (repo *GormRepository) CreateUser(user *models.User) (*models.User, error) {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if exist, err := gormutil.RecordExists(tx, &user); err != nil {
			return err
		} else if exist {
			return ErrAlreadyExists
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	repo.bus.Publish(event.UserCreated, event.UserCreatedEvent{
		UserID:   user.ID,
		Username: user.Username,
	})
	return user, nil
}

// GetUser : Get the user with this id
func (repo *GormRepository) GetUser(userID uint) (*models.User, error) {
	user := models.User{}
	user.ID = userID
	repo.db.Find(&user)
	return &user, nil
}

// GetDevicesByUserID : Get all devices of the user
func (repo *GormRepository) GetDevicesByUserID(userID uint) ([]models.Device, error) {
	user := models.User{}
	user.ID = userID
	repo.db.Find(&user)
	return user.Devices, nil
}

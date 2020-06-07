package repository

import (
	"github.com/cynt4k/wygops/internal/event"
	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/pkg/util/gormutil"
	"github.com/cynt4k/wygops/pkg/util/structs"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
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

	msg, _ := structs.StructToMap(event.UserCreatedEvent{
		UserID:   user.ID,
		Username: user.Username,
	})

	repo.hub.Publish(hub.Message{
		Name:   event.UserCreated,
		Fields: msg,
	})
	return user, nil
}

// GetUser : Get the user with this id
func (repo *GormRepository) GetUser(userID uint) (*models.User, error) {
	user := models.User{
		ID: userID,
	}
	// if err := repo.db.First(&user).Error; err != nil {
	// 	return nil, err
	// }
	// return &user, nil
	return getUser(repo.db, false, &user)
}

// GetUserWithDevices : Get the user with his devices
func (repo *GormRepository) GetUserWithDevices(userID uint) (*models.User, error) {
	user := models.User{
		ID: userID,
	}
	return getUser(repo.db, true, &user)

}

func getUser(tx *gorm.DB, withDevices bool, where ...interface{}) (*models.User, error) {
	var user models.User
	if withDevices {
		tx = tx.Preload("Devices")
	}
	if err := tx.First(&user, where...).Error; err != nil {
		return nil, convertError(err)
	}
	return &user, nil
}

// DeleteUser : Delete an user
func (repo *GormRepository) DeleteUser(userID uint) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Where(&models.UserGroup{UserID: userID}).Delete(&models.UserGroup{}).Error; err != nil {
			return err
		}
		if err := tx.Where(&models.Device{UserID: userID}).Delete(&models.Device{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
	if err != nil {
		return err
	}

	msg, _ := structs.StructToMap(event.UserDeletedEvent{
		UserID: userID,
	})

	repo.hub.Publish(hub.Message{
		Name:   event.UserDeleted,
		Fields: msg,
	})
	return nil
}

package repository

import (
	"github.com/cynt4k/wygops/internal/event"
	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/pkg/util/gormutil"
	"github.com/cynt4k/wygops/pkg/util/structs"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
)

// CreateDevice : Create an device
func (repo *GormRepository) CreateDevice(device *models.Device) (*models.Device, error) {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if exist, err := gormutil.RecordExists(tx, &device); err != nil {
			return err
		} else if exist {
			return ErrAlreadyExists
		}

		if err := tx.Create(&device).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	msg, _ := structs.StructToMap(event.DeviceCreatedEvent{
		DeviceID: device.ID,
		UserID:   device.UserID,
	})

	repo.hub.Publish(hub.Message{
		Name:   event.DeviceCreated,
		Fields: msg,
	})

	return device, nil
}

// GetDevices : Get all devices
func (repo *GormRepository) GetDevices() ([]*models.Device, error) {
	var devices []*models.Device

	if err := repo.db.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// FindDevices : Find all devices by pattern
func (repo *GormRepository) FindDevices(where ...interface{}) ([]*models.Device, error) {
	var devices []*models.Device
	if err := repo.db.Find(&devices, where...).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

// GetDevicesByUserID : Get all devices of the user
func (repo *GormRepository) GetDevicesByUserID(userID uint) ([]*models.Device, error) {
	user := models.User{
		ID: userID,
	}
	if err := repo.db.First(&user).Error; err != nil {
		return nil, err
	}
	return user.Devices, nil
}

// DeleteDevice : Delete a device
func (repo *GormRepository) DeleteDevice(deviceID uint) error {
	var device models.Device
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&device, &models.Device{ID: deviceID}).Error; err != nil {
			return err
		}
		result := tx.Delete(&models.Device{}, &models.Device{ID: deviceID})
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

	msg, _ := structs.StructToMap(event.DeviceDeletedEvent{
		DeviceID: deviceID,
		UserID:   device.UserID,
	})

	repo.hub.Publish(hub.Message{
		Name:   event.DeviceDeleted,
		Fields: msg,
	})

	return nil
}

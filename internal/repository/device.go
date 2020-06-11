package repository

import "github.com/cynt4k/wygops/internal/models"

// DeviceRepository : Device repository to predfine the interfaces
type DeviceRepository interface {
	CreateDevice(*models.Device) (*models.Device, error)
	// UpdateDevice(*models.Device) (*models.Device, error)
	GetDevices() ([]*models.Device, error)
	GetDevicesByUserID(uint) ([]*models.Device, error)
	DeleteDevice(deviceID uint) error
}

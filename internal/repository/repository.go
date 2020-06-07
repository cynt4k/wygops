package repository

// Repository : Repository Interface to be implemented
type Repository interface {
	Sync() (bool, error)
	UserRepository
	GroupRepository
	DeviceRepository
}

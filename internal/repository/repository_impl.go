package repository

import (
	evbus "github.com/asaskevich/EventBus"
	"github.com/jinzhu/gorm"
)

// GormRepository : Repository struct to handle
type GormRepository struct {
	db  *gorm.DB
	bus evbus.Bus
}

// NewGormRepository : Create a new repository to handle orm operations
func NewGormRepository(db *gorm.DB, bus *evbus.Bus) (Repository, error) {
	repo := &GormRepository{
		db:  db,
		bus: *bus,
	}

	return repo, nil
}

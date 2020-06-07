package repository

import (
	evbus "github.com/asaskevich/EventBus"
	"github.com/cynt4k/wygops/internal/migration"
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

// Sync : Synchronice database with migration and default settings
func (repo *GormRepository) Sync() (init bool, err error) {
	if err := migration.Migrate(repo.db); err != nil {
		return false, err
	}
	return true, nil
}

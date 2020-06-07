package repository

import (
	"github.com/cynt4k/wygops/internal/migration"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
)

// GormRepository : Repository struct to handle
type GormRepository struct {
	db  *gorm.DB
	hub *hub.Hub
}

// NewGormRepository : Create a new repository to handle orm operations
func NewGormRepository(db *gorm.DB, hub *hub.Hub) (Repository, error) {
	repo := &GormRepository{
		db:  db,
		hub: hub,
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

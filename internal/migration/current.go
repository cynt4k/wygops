package migration

import (
	"github.com/cynt4k/wygops/internal/models"
	"gopkg.in/gormigrate.v1"
)

// Migrations : Return the migrations to be executed
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{}
}

// AllTables : All current used tables
func AllTables() []interface{} {
	return []interface{}{
		&models.Device{},
		&models.UserGroup{},
		&models.User{},
		&models.Group{},
	}
}

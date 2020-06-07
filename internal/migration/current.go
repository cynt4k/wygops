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

// AllForeignKeys : All current foreign keys
func AllForeignKeys() [][5]string {
	return [][5]string{
		// Table, Key, Reference, OnDelete, OnUpdate
		{"devices", "user_id", "user(id)", "CASCADE", "CASCADE"},
	}
}

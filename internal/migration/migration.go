package migration

import (
	"github.com/jinzhu/gorm"
	gomigrate "gopkg.in/gormigrate.v1"
)

// Migrate : Migrate database version
func Migrate(db *gorm.DB) error {
	m := gomigrate.New(db, &gomigrate.Options{
		TableName:                 "migration",
		IDColumnName:              "id",
		IDColumnSize:              190,
		UseTransaction:            false,
		ValidateUnknownMigrations: true,
	}, Migrations())

	m.InitSchema(func(db *gorm.DB) error {
		if err := db.AutoMigrate(AllTables()...).Error; err != nil {
			return err
		}

		return nil
	})

	return m.Migrate()
}

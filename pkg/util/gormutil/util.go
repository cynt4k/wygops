package gormutil

import "github.com/jinzhu/gorm"

// RecordExists : Check if the provided record does already exist
func RecordExists(db *gorm.DB, where interface{}, tableName ...string) (exists bool, err error) {
	if len(tableName) > 0 {
		db = db.Table(tableName[0])
	} else {
		db = db.Model(where)
	}
	return Exists(db.Where(where))
}

// Exists : Check if matches more than zero
func Exists(db *gorm.DB) (exists bool, err error) {
	n, err := Count(db.Limit(1))
	return n > 0, err
}

// Count : Count number of query results
func Count(db *gorm.DB) (n int, err error) {
	return n, db.Count(&n).Error
}

package mysql

import (
	"fmt"

	"github.com/cynt4k/wygops/cmd"
	"github.com/cynt4k/wygops/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type mysql struct{}

var (
	connection = mysql{}
)

// Init : Initialize the mysql connection
func Init() (*gorm.DB, error) {
	return connection.init()
}

func (c *mysql) init() (*gorm.DB, error) {
	config := cmd.GetConfig().Database
	// user:password@/dbname?charset=utf8&parseTime=True&loc=Local
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database))

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Device{}).Error; err != nil {
		return nil, err
	}

	return db, err
}

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/cynt4k/wygops/cmd/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func readDefaultConfig(configDir string, configFile *config.Config) error {
	yamlFile, err := ioutil.ReadFile(filepath.Join(configDir, "default.yaml"))

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &configFile)

	if err != nil {
		return err
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(yamlFile))

	if err != nil {
		return err
	}

	return nil
}

func getDatabase(databaseConfig *config.Database, mode bool) (*gorm.DB, error) {
	engine, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true",
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Database,
	))
	if err != nil {
		return nil, err
	}

	engine.BlockGlobalUpdate(true)
	engine.LogMode(mode)
	return engine.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"), nil
}

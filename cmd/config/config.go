package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	configFile string
	c          *Config
)

// GeneralUser : Config struct for user settings
type GeneralUser struct {
	MaxDevices int8 `yaml:"maxDevices"`
}

// General : Config struct for general settings
type General struct {
	User GeneralUser `yaml:"user"`
}

// Database : Config struct for the database
type Database struct {
	Host     string `yaml:"host"`
	Port     int16  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// ProviderLdap : Config struct for the LDAP provider
type ProviderLdap struct {
	Enabled      bool     `yaml:"enabled"`
	Type         string   `yaml:"type"`
	Host         string   `yaml:"host"`
	Port         int16    `yaml:"port"`
	BindDn       string   `yaml:"bindDn"`
	BindPassword string   `yaml:"bindPassword"`
	BaseDn       string   `yaml:"baseDn"`
	GroupFilter  string   `yaml:"groupFilter"`
	UserFilter   string   `yaml:"userFilter"`
	UserAttr     []string `yaml:"userAttr"`
	GroupAttr    []string `yaml:"groupAttr"`
	UserRDN      string   `yaml:"userRdn"`
	GroupRDN     string   `yaml:"groupRdn"`
}

// Provider : Config struct for the Providers
type Provider struct {
	Ldap ProviderLdap `yaml:"ldap"`
}

// API : Config struct for the api
type API struct {
	Host string `yaml:"host"`
	Port int16  `yaml:"port"`
}

// Config type for the config file to be handeld
type Config struct {
	DevMode  bool     `yaml:"dev"`
	General  General  `yaml:"general"`
	Database Database `yaml:"database"`
	Provider Provider `yaml:"provider"`
	API      API      `yaml:"api"`
}

var (
	config = Config{}
)

// GetConfig : Return the viper config
func GetConfig() *Config {
	return &config
}

func readDefaultConfig(configDir string) error {
	yamlFile, err := ioutil.ReadFile(filepath.Join(configDir, "default.yaml"))

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &config)

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

func (c Config) getDatabase() (*gorm.DB, error) {
	engine, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	))
	if err != nil {
		return nil, err
	}

	engine.BlockGlobalUpdate(true)
	engine.LogMode(c.DevMode)
	return engine.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"), nil
}

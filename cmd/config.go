package cmd

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// ConfigGeneralUser : Config struct for user settings
type ConfigGeneralUser struct {
	MaxDevices int8 `yaml:"maxDevices"`
}

// ConfigGeneral : Config struct for general settings
type ConfigGeneral struct {
	User ConfigGeneralUser `yaml:"user"`
}

// ConfigDatabase : Config struct for the database
type ConfigDatabase struct {
	Host     string `yaml:"host"`
	Port     int16  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// ConfigProviderLdap : Config struct for the LDAP provider
type ConfigProviderLdap struct {
	Enabled      bool   `yaml:"enabled"`
	URI          string `yaml:"uri"`
	BindDn       string `yaml:"bindDn"`
	BindPassword string `yaml:"bindPassword"`
	BaseDn       string `yaml:"baseDn"`
	GroupFilter  string `yaml:"groupFilter"`
}

// ConfigProvider : Config struct for the Providers
type ConfigProvider struct {
	Ldap ConfigProviderLdap `yaml:"ldap"`
}

// Config type for the config file to be handeld
type Config struct {
	General  ConfigGeneral  `yaml:"general"`
	Database ConfigDatabase `yaml:"database"`
	Provider ConfigProvider `yaml:"provider"`
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

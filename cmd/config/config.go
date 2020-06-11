package config

import (
	_ "github.com/go-sql-driver/mysql"
)

var (
	configFile string
	c          *Config
)

// GeneralUser : Config struct for user settings
type GeneralUser struct {
	MaxDevices int8 `yaml:"maxDevices"`
}

// GeneralSubnet : Config struct for the subnet
type GeneralSubnet struct {
	V4        string `yaml:"v4"`
	V6        string `yaml:"v6"`
	GatewayV4 string `yaml:"gatewayV4"`
	GatewayV6 string `yaml:"gatewayV6"`
}

// General : Config struct for general settings
type General struct {
	User   GeneralUser   `yaml:"user"`
	Subnet GeneralSubnet `yaml:"subnet"`
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

// Wireguard : Config struct for wireguard
type Wireguard struct {
	Interface string `yaml:"interface"`
}

// Config type for the config file to be handeld
type Config struct {
	DevMode   bool      `yaml:"dev"`
	General   General   `yaml:"general"`
	Wireguard Wireguard `yaml:"wireguard"`
	Database  Database  `yaml:"database"`
	Provider  Provider  `yaml:"provider"`
	API       API       `yaml:"api"`
}

var (
	config = Config{}
)

// GetConfig : Return the viper config
func GetConfig() *Config {
	return &config
}

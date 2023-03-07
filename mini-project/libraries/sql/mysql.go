package sql

import (
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Host     string `mapstructure:"MYSQL_HOST"`
	Port     int    `mapstructure:"MYSQL_PORT"`
	Username string `mapstructure:"MYSQL_USERNAME"`
	Password string `mapstructure:"MYSQL_PASSWORD"`
	Schema   string `mapstructure:"MYSQL_SCHEMA"`
}

func NewMySQLConn(cfg *MySQLConfig) (db *gorm.DB, err error) {
	return nil, nil
}

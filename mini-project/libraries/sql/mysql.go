package sql

import (
	"fmt"

	"gorm.io/driver/mysql"
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
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Schema,
	)

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
	}

	if db, err = gorm.Open(mysql.Open(dsn), gormConfig); err != nil {
		return
	}

	return
}

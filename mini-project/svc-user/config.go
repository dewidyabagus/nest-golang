package main

import (
	"log"

	"nest-golang/libraries/sql"

	"github.com/spf13/viper"
)

type AppConfig struct {
	AppName      string `mapstructure:"APP_NAME"`
	ListenPort   int    `mapstructure:"APP_PORT"`
	Environment  string `mapstructure:"ENV"`
	JwtSecretKey string `mapstructure:"JWT_SECRET_KEY"`
	Database     sql.MySQLConfig
	Notify       NotifyConfig
}

type NotifyConfig struct {
	Host   string `mapstructure:"NOTIFY_HOST"`
	ApiKey string `mapstructure:"NOTIFY_API_KEY"`
}

const defaultPort = 7001
const defaultPath = "./svc-user"
const defaultEnv = "development"

func LoadConfig() *AppConfig {
	viper.AddConfigPath(".")
	viper.AddConfigPath(defaultPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("read config failed, error:", err.Error())
	}

	viper.SetDefault("ENV", defaultEnv)
	viper.SetDefault("APP_PORT", defaultPort)

	config := &AppConfig{}
	viper.Unmarshal(config)
	viper.Unmarshal(&config.Database)
	viper.Unmarshal(&config.Notify)

	return config
}

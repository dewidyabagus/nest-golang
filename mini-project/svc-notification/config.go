package main

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	AppName     string `mapstructure:"APP_NAME"`
	ListenPort  int    `mapstructure:"APP_PORT"`
	Environment string `mapstructure:"ENV"`
	SecretKey   string `mapstructure:"SECRET_KEY"`
	MailConfig  MailConfig
}

type MailConfig struct {
	SMTPHost   string `mapstructure:"EMAIL_SMTP_HOST"`
	SMTPPort   int    `mapstructure:"EMAIL_SMTP_PORT"`
	SenderName string `mapstructure:"EMAIL_SENDER_NAME"`
	Email      string `mapstructure:"EMAIL_AUTH_EMAIL"`
	Password   string `mapstructure:"EMAIL_AUTH_PASSWORD"`
}

const defaultPort = 7002
const defaultPath = "./svc-notification"
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
	viper.Unmarshal(&config.MailConfig)

	return config
}

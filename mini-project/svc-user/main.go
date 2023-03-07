package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nest-golang/libraries/sql"
	"nest-golang/libraries/validator"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := LoadConfig()

	db, err := sql.NewMySQLConn(&cfg.Database)
	if err != nil {
		log.Fatalln("new mysql conn failed, error:", err.Error())
	}
	instDB, _ := db.DB()
	defer instDB.Close()

	// Add table suffix when creating tables
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

	notifyAggregator := NewNotify(cfg.Notify)
	userRepository := NewRepository(db)
	userController := NewHandler(cfg.JwtSecretKey, userRepository, notifyAggregator)

	e := echo.New()
	e.Validator = validator.New(validator.Options{TagNameFunc: "json"})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello User!")
	})
	e.GET("/health", func(c echo.Context) error {
		var info struct {
			Service  string `json:"service"`
			Time     string `json:"time"`
			Database string `json:"database"`
		}
		info.Service = cfg.AppName
		info.Time = time.Now().Format("2006-01-02 15:04:05")

		if err := instDB.Ping(); err != nil {
			info.Database = "down"
			return c.JSON(http.StatusInternalServerError, info)
		}

		info.Database = "up"
		return c.JSON(http.StatusOK, info)
	})
	e.POST("/register", userController.PostUser)
	e.POST("/login", userController.Login)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.ListenPort)); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalln("Shutting down server failed, error:", err.Error())
	}
}

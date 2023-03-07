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

	"nest-golang/libraries/validator"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const badRequest = "BAD_REQUEST"
const invalidData = "INVALID_DATA"
const serverError = "INTERNAL_SERVER_ERROR"

func main() {
	cfg := LoadConfig()

	mailer := NewMailer(cfg.MailConfig)
	handler := NewHandler(mailer)

	e := echo.New()
	e.Validator = validator.New(validator.Options{TagNameFunc: "json"})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	route := e.Group("email", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-API-KEY",
		Validator: func(auth string, c echo.Context) (bool, error) {
			return cfg.SecretKey == auth, nil
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Notification!")
	})

	e.GET("/health", func(c echo.Context) error {
		var info struct {
			Service string `json:"service"`
			Time    string `json:"time"`
			Mailer  string `json:"mailer"`
		}
		info.Service = cfg.AppName
		info.Time = time.Now().Format("2006-01-02 15:04:05")

		if err := mailer.Ping(); err != nil {
			info.Mailer = "down"
			return c.JSON(http.StatusInternalServerError, info)
		}

		info.Mailer = "up"
		return c.JSON(http.StatusOK, info)
	})

	route.POST("/:notif_type", handler.EmailNotify)

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

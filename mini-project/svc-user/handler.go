package main

import (
	"context"
	"net/http"
	"time"

	"nest-golang/libraries/validator"

	"github.com/golang-jwt/jwt"
	echo "github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type servicer interface {
	PostUser(ctx context.Context, u User) error
	GetUserWithEmail(ctx context.Context, email string) (user User, err error)
}

type notifier interface {
	EmailNotify(notif string, info PayloadNotify) error
}

const badRequest = "BAD_REQUEST"
const invalidData = "INVALID_DATA"
const serverError = "INTERNAL_SERVER_ERROR"

type handler struct {
	jwtSecretKey string
	service      servicer
	notify       notifier
}

func NewHandler(jwtKey string, s servicer, n notifier) *handler {
	return &handler{service: s, notify: n, jwtSecretKey: jwtKey}
}

func (h *handler) PostUser(c echo.Context) error {
	user := new(InsertUser)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: badRequest, Errors: err.Error()})
	}

	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: invalidData, Errors: validator.ErrorFormTranslator(err)})
	}

	if err := h.service.PostUser(c.Request().Context(), user.RegisterUser()); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: serverError, Errors: err.Error()})
	}

	go h.notify.EmailNotify("register", PayloadNotify{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})

	return c.JSON(http.StatusCreated, Response{Message: "sukses menambahkan user"})
}

func (h *handler) Login(c echo.Context) error {
	user := new(UserLogin)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: badRequest, Errors: err.Error()})
	}

	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: invalidData, Errors: validator.ErrorFormTranslator(err)})
	}

	res, err := h.service.GetUserWithEmail(c.Request().Context(), user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: serverError, Errors: err.Error()})

	} else if res.ID == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: badRequest, Errors: "user dan password tidak valid"})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: badRequest, Errors: "password tidak valid"})
	}

	eJWT := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  res.ID,
			"exp": time.Now().Add(time.Hour).Unix(),
		},
	)

	token, _ := eJWT.SignedString([]byte(h.jwtSecretKey))

	go h.notify.EmailNotify("login", PayloadNotify{
		Email:     user.Email,
		FirstName: res.FirstName,
		LastName:  res.LastName,
	})

	return c.JSON(http.StatusOK, Response{Message: "berhasil login", Data: echo.Map{"token": token}})
}

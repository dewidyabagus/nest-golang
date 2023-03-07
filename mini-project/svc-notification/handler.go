package main

import (
	"nest-golang/libraries/validator"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

type handler struct {
	mailer *mailer
}

func NewHandler(ml *mailer) *handler {
	return &handler{ml}
}

func (h *handler) EmailNotify(c echo.Context) (err error) {
	notifType := c.Param("notif_type")
	if notifType != "login" && notifType != "register" {
		return c.JSON(http.StatusBadRequest, Response{Status: badRequest, Errors: "tipe notifikasi tidak terdaftar"})
	}

	info := UserNotify{}
	if err := c.Bind(&info); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: badRequest, Errors: err.Error()})
	}

	if err := c.Validate(info); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Status: invalidData, Errors: validator.ErrorFormTranslator(err)})
	}

	if notifType == "login" {
		err = h.mailer.SendLoginNotify(info)
	} else {
		err = h.mailer.SendAccountActivation(info)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Status: serverError, Errors: err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "sukses mengirimkan notifikasi"})
}

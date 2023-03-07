package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Pengaliasan type dari package github.com/go-playground/validator/v10
type ValidationErrors = validator.ValidationErrors

type Options struct {
	TagNameFunc string // Nama tag untuk diambil nilainya
}

type validation struct {
	validator *validator.Validate
}

func New(opts Options) *validation {
	vld := validator.New()

	if tag := strings.TrimSpace(opts.TagNameFunc); tag != "" {
		vld.RegisterTagNameFunc(func(field reflect.StructField) string {
			if name := strings.TrimSpace(field.Tag.Get(tag)); name != "" {
				return name
			}
			return field.Name // default
		})
	}

	return &validation{validator: vld}
}

func (v *validation) Validate(s interface{}) error {
	return v.validator.Struct(s)
}

func ErrorFormTranslator(err error) interface{} {
	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err

	} else if len(valErrors) == 0 {
		return nil
	}

	fields := map[string]string{}
	for _, err := range valErrors {
		switch err.Tag() {
		default:
			fields[err.Field()] = err.Error()

		case "required":
			fields[err.Field()] = "is required"

		case "email":
			fields[err.Field()] = "is not valid email"

		case "gte":
			fields[err.Field()] = fmt.Sprintf("value must be greater than equal %s", err.Param())

		case "lte":
			fields[err.Field()] = fmt.Sprintf("value must be lower than equal %s", err.Param())

		case "eqfield":
			fields[err.Field()] = fmt.Sprintf("value must be same as field %s", err.Param())
		}
	}

	return fields
}

package model

import (
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

var (
	CustomValidators = map[string]customValidator{
		"enum": {
			Message:   "This field should be in: %v",
			Validator: validateEnum,
		},
		"required": {
			Message: "This field is required and cannot be empty",
		},
	}
)

type customValidator struct {
	Message   string
	Validator validator.Func
}

func validateEnum(fl validator.FieldLevel) bool {
	for _, v := range strings.Split(fl.Param(), " ") {
		if v == fl.Field().String() {
			return true
		}
	}
	return false
}

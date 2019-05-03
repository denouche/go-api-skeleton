package validators

import (
	"context"
	"regexp"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"gopkg.in/go-playground/validator.v9"
)

var (
	ContextKeyValidator = "ContextKeyValidator"

	CustomValidators = map[string]customValidator{
		"regexp": {
			Message:   "This field should match the following pattern: %v",
			Validator: validateRegexp,
		},
		"required": {
			Message: "This field is required and cannot be empty",
		},
	}
)

type ValidationContext struct {
	DB dao.Database
}

type customValidator struct {
	Message   string
	Validator validator.FuncCtx
}

func validateRegexp(ctx context.Context, fl validator.FieldLevel) bool {
	return regexp.MustCompile(fl.Param()).MatchString(fl.Field().String())
}

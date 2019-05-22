package validators

import (
	"context"
	"reflect"
	"regexp"
	"strings"

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

type validationContext struct {
	DB dao.Database
}

type customValidator struct {
	Message   string
	Validator validator.FuncCtx
}

func validateRegexp(ctx context.Context, fl validator.FieldLevel) bool {
	return regexp.MustCompile(fl.Param()).MatchString(fl.Field().String())
}

func NewValidator() *validator.Validate {
	va := validator.New()

	va.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)
		if len(name) < 1 {
			return ""
		}
		return name[0]
	})

	for k, v := range CustomValidators {
		if v.Validator != nil {
			va.RegisterValidationCtx(k, v.Validator)
		}
	}

	return va
}

func NewContextWithValidationContext(parentCtx context.Context, db dao.Database) context.Context {
	vc := &validationContext{
		DB: db,
	}
	return context.WithValue(parentCtx, ContextKeyValidator, vc)
}

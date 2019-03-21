package validators

import (
	"context"
	"strings"

	"github.com/denouche/go-api-skeleton/storage/dao"
	"gopkg.in/go-playground/validator.v9"
)

var (
	ContextKeyValidator = "ContextKeyValidator"

	CustomValidators = map[string]customValidator{
		"enum": {
			Message:   "This field should be in: %v",
			Validator: validateEnum,
		},
		"required": {
			Message: "This field is required and cannot be empty",
		},
		//"something-exists": {
		//	Message:   "This field should exists",
		//	Validator: validateSomethingExists,
		//},
	}
)

type ValidationContext struct {
	DB dao.Database
}

type customValidator struct {
	Message   string
	Validator validator.FuncCtx
}

func validateEnum(ctx context.Context, fl validator.FieldLevel) bool {
	for _, v := range strings.Split(fl.Param(), " ") {
		if v == fl.Field().String() {
			return true
		}
	}
	return false
}

//func validateSomethingExists(ctx context.Context, fl validator.FieldLevel) bool {
//	vc := ctx.Value(ContextKeyValidator)
//	if vc == nil {
//		return false
//	}
//
//	validationContext, ok := vc.(*ValidationContext)
//	if !ok {
//		return false
//	}
//
//	// here in the validationContext you have the DB for any db required validation
//}

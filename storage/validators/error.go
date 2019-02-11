package validators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/denouche/go-api-skeleton/storage/model"
	"github.com/denouche/go-api-skeleton/utils"
	"gopkg.in/go-playground/validator.v9"
)

var regexpValidatorNamespacePrefix = regexp.MustCompile(`^\w+\.`)

func NewDataValidationAPIError(err error) model.APIError {
	apiErr := model.ErrDataValidation
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			utils.GetLogger(nil).Errorw("InvalidValidationError",
				"templateAPIErr", apiErr)
		} else {
			for _, e := range err.(validator.ValidationErrors) {
				reason := e.Tag()
				if _, ok := CustomValidators[e.Tag()]; ok {
					reason = truncatingSprintf(CustomValidators[e.Tag()].Message, e.Param())
				}

				namespaceWithoutStructName := regexpValidatorNamespacePrefix.ReplaceAllString(e.Namespace(), "")
				fe := model.FieldError{
					Field:       namespaceWithoutStructName,
					Constraint:  e.Tag(),
					Description: reason,
				}
				apiErr.Details = append(apiErr.Details, fe)
			}
		}
	}
	return apiErr
}

// truncatingSprintf is used as fmt.Sprintf but allow to truncate the additional parameters given when there is more parameters than %v in str
func truncatingSprintf(str string, args ...interface{}) string {
	n := strings.Count(str, "%v")
	return fmt.Sprintf(str, args[:n]...)
}

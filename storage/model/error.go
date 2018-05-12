package model

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/denouche/go-api-skeleton/storage/validators"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

var regexpValidatorNamespacePrefix = regexp.MustCompile(`^\w+\.`)

var (
	// 400
	ErrBadRequestFormat = APIError{
		Type:        "bad_format",
		HTTPCode:    http.StatusBadRequest,
		Description: "unable to read request body, please check that the json is valid",
	}
	ErrDataValidation = APIError{
		Type:        "data_validation",
		HTTPCode:    http.StatusBadRequest,
		Description: "the data are not valid",
	}

	// 404
	ErrNotFound = APIError{
		Type:     "not_found",
		HTTPCode: http.StatusNotFound,
	}

	// 40x
	ErrAlreadyExists = APIError{
		Type:     "already_exists",
		HTTPCode: http.StatusConflict,
	}

	// 50x
	ErrInternalServer = APIError{
		Type:     "internal_server_error",
		HTTPCode: http.StatusInternalServerError,
	}
)

type APIError struct {
	HTTPCode    int                 `json:"-"`
	Type        string              `json:"error"`
	Description string              `json:"error_description"`
	Details     []FieldError        `json:"error_details,omitempty"`
	Headers     map[string][]string `json:"-"`
}

func NewDataValidationAPIError(err error) APIError {
	apiErr := ErrDataValidation
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logrus.WithError(err).WithField("templateAPIErr", apiErr).Error("InvalidValidationError")
		} else {
			for _, e := range err.(validator.ValidationErrors) {
				reason := e.Tag()
				if _, ok := validators.CustomValidators[e.Tag()]; ok {
					reason = truncatingSprintf(validators.CustomValidators[e.Tag()].Message, e.Param())
				}

				namespaceWithoutStructName := regexpValidatorNamespacePrefix.ReplaceAllString(e.Namespace(), "")
				fe := FieldError{
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

type FieldError struct {
	Field       string `json:"field"`
	Constraint  string `json:"constraint"`
	Description string `json:"description"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("error : %d, %s, %s, %v", e.HTTPCode, e.Type, e.Description, e.Details)
}

// truncatingSprintf is used as fmt.Sprintf but allow to truncate the additional parameters given when there is more parameters than %v in str
func truncatingSprintf(str string, args ...interface{}) string {
	n := strings.Count(str, "%v")
	return fmt.Sprintf(str, args[:n]...)
}

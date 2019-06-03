package model

import (
	"fmt"
	"net/http"
)

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
	ErrVersionMismatched = APIError{
		Type:        "precondition_failed",
		HTTPCode:    http.StatusPreconditionFailed,
		Description: "Model version mismatched",
	}

	// 50x
	ErrInternalServer = APIError{
		Type:     "internal_server_error",
		HTTPCode: http.StatusInternalServerError,
	}
)

// @openapi:schema
type APIError struct {
	HTTPCode    int                 `json:"-"`
	Type        string              `json:"error"`
	Description string              `json:"error_description"`
	Details     []FieldError        `json:"error_details,omitempty"`
	Headers     map[string][]string `json:"-"`
}

// @openapi:schema
type FieldError struct {
	Field       string `json:"field"`
	Constraint  string `json:"constraint"`
	Description string `json:"description"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("error : %d, %s, %s, %v", e.HTTPCode, e.Type, e.Description, e.Details)
}

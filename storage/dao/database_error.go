package dao

import (
	"fmt"
)

type Type int

const (
	ErrTypeNotFound Type = iota
	ErrTypeDuplicate
	ErrTypeForeignKeyViolation
)

type DAOError struct {
	Cause error
	Type  Type
}

func newDAOError(t Type, cause error) error {
	return &DAOError{
		Type:  t,
		Cause: cause,
	}
}

func (e *DAOError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Type %d: %s", e.Type, e.Cause.Error())
	}
	return fmt.Sprintf("Type %d: no cause given", e.Type)
}

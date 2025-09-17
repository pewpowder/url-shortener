package errors

import (
	"errors"

	"gorm.io/gorm"
)

func GormErrorToErrorType(err error) ErrorType {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrTypeNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrTypeConflict
	case errors.Is(err, gorm.ErrForeignKeyViolated) || errors.Is(err, gorm.ErrCheckConstraintViolated):
		return ErrTypeValidation
	case errors.Is(err, gorm.ErrPrimaryKeyRequired) ||
		errors.Is(err, gorm.ErrModelValueRequired) ||
		errors.Is(err, gorm.ErrModelAccessibleFieldsRequired) ||
		errors.Is(err, gorm.ErrSubQueryRequired) ||
		errors.Is(err, gorm.ErrInvalidData) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrInvalidValue) ||
		errors.Is(err, gorm.ErrInvalidValueOfLength):
		return ErrTypeValidation
	default:
		return ErrTypeInternal
	}
}

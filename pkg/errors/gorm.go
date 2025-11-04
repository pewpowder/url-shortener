package errors

import (
	"errors"

	"gorm.io/gorm"
)

func MapGormErrToErrCodeWithMessage(err error) (ErrorCode, string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrCodeNotFound, "RESOURCE_NOT_FOUND"
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrCodeConflict, "RESOURCE_ALREADY_EXISTS"
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return ErrCodeBadRequest, "INVALID_REFERENCE_TO_ANOTHER_RESOURCE"
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		return ErrCodeBadRequest, "DATA_VALIDATION_FAILED"
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return ErrCodeBadRequest, "UNSUPPORTED_RELATION"
	case errors.Is(err, gorm.ErrMissingWhereClause) ||
		errors.Is(err, gorm.ErrPrimaryKeyRequired) ||
		errors.Is(err, gorm.ErrModelValueRequired) ||
		errors.Is(err, gorm.ErrModelAccessibleFieldsRequired) ||
		errors.Is(err, gorm.ErrSubQueryRequired) ||
		errors.Is(err, gorm.ErrInvalidData) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrInvalidValue) ||
		errors.Is(err, gorm.ErrInvalidValueOfLength) ||
		errors.Is(err, gorm.ErrEmptySlice) ||
		errors.Is(err, gorm.ErrPreloadNotAllowed):
		return ErrCodeBadRequest, "INVALID_REQUEST"
	default:
		return ErrCodeInternal, "INTERNAL_ERROR"
	}
}

func NewServiceErrorFromGorm(message string, err error) error {
	code, codeMessage := MapGormErrToErrCodeWithMessage(err)

	return NewServiceError(message, code, codeMessage, err)
}

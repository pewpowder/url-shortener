package errors

import "gorm.io/gorm"

func GormErrorToErrorType(gormErr error) ErrorType {
	if gormErr == nil {
		return ErrTypeInternal
	}

	switch gormErr {
	case gorm.ErrRecordNotFound:
		return ErrTypeNotFound
	case gorm.ErrDuplicatedKey:
		return ErrTypeConflict
	case gorm.ErrForeignKeyViolated,
		gorm.ErrCheckConstraintViolated:
		return ErrTypeValidation
	case
		gorm.ErrPrimaryKeyRequired, gorm.ErrModelValueRequired,
		gorm.ErrModelAccessibleFieldsRequired, gorm.ErrSubQueryRequired,
		gorm.ErrInvalidData, gorm.ErrInvalidField, gorm.ErrInvalidValue,
		gorm.ErrInvalidValueOfLength:
		return ErrTypeValidation
	default:
		return ErrTypeInternal
	}
}

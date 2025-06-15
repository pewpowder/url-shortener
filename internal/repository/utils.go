package repository

import (
	"errors"

	"gorm.io/gorm"
)

func IsGormNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

package entity

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name  string `gorm:"size:50;index;unique"`
	Color string `gorm:"size:7;"`
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

// TODO: change implementation to count all visits not per day
type LinkStats struct {
	gorm.Model
	LinkID     uint      `json:"link_id"`
	Date       time.Time `gorm:"index"`
	ClickCount int
	Countries  map[string]int // country name as key, count as value (maybe serialize to JSON?)
}

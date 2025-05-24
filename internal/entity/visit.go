package entity

import (
	"time"

	"gorm.io/gorm"
)

type Visit struct {
	gorm.Model
	LinkID    uint      `json:"link_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

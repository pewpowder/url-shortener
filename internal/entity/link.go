package entity

import (
	"time"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	ShortCode   string      `json:"short_code" gorm:"uniqueIndex"`
	OriginalURL string      `json:"original_url"`
	ExpiresAt   time.Time   `json:"expires_at,omitempty"`
	IsActive    bool        `json:"is_active"`
	Stats       []LinkStats `json:"stats" gorm:"type:jsonb;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Visits      []Visit     `json:"visit" gorm:"type:jsonb;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// UserID      uint        `json:"user_id"`
}

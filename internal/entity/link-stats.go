package entity

import "gorm.io/gorm"

type LinkStats struct {
	gorm.Model
	LinkID      uint   `json:"link_id"`
	DailyVisits []byte `json:"daily_clicks" gorm:"type:jsonb"`
}

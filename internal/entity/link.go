package entity

import (
	"time"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	ShortCode    string     `gorm:"uniqueIndex;size:10;check:short_code <> ''"`
	OriginalURL  string     `gorm:"not null;check:original_url <> ''"`
	ExpiresAt    *time.Time `gorm:"index"`
	IsActive     bool       `gorm:"default:true"`
	IsPrivate    bool       `gorm:"default:false"`
	UserID       uint       // Integration with Auth-Service
	PasswordHash *string    `gorm:"size:64"` // PasswordHash for private links
	MaxClicks    int        `gorm:"default:0"`
	Tags         []Tag      `gorm:"many2many:link_tags;"`
}

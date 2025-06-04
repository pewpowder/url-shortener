package entity

import "time"

type User struct {
	ID         uint   `json:"id"` // Foreign key (from auth-service)
	Email      string `json:"email"`
	LastUpdate time.Time
}

package entity

type User struct {
	ID    uint   `json:"id"` // Foreign key (from auth-service)
	Email string `json:"email"`
}

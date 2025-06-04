package entity

type Tag struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"uniqueIndex"`
	Color string `gorm:"size:7;default:'#4f46e5'"`
}

// TODO: Add additional endpoint for managing tags (when craete link we just pass tag name and check if it exists)

package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type LinkStats struct {
	gorm.Model
	LinkID     uint `gorm:"index" json:"link_id"`
	ClickCount int
	Countries  Countries `gorm:"type:jsonb"`
}

type Countries []Country

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (c *Countries) Scan(value any) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("неподдерживаемый тип: %T", value)
	}

	return json.Unmarshal(data, c)
}

func (c Countries) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}

	return json.Marshal(c)
}

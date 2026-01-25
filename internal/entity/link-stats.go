package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

var LINK_STATS_SCHEMA = `
	CREATE TABLE IF NOT EXISTS link_stats (
		id SERIAL PRIMARY KEY,
		link_id INTEGER NOT NULL,
		click_count BIGINT NOT NULL DEFAULT 0,
		countries JSONB NOT NULL DEFAULT '[]'::jsonb,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ,
		deleted_at TIMESTAMPTZ,
		CONSTRAINT fk_link_stats_link FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
	);
	
	CREATE INDEX IF NOT EXISTS idx_link_stats_link_id ON link_stats(link_id);
	CREATE INDEX IF NOT EXISTS idx_link_stats_deleted_at ON link_stats(deleted_at);
`

type LinkStats struct {
	ID         uint       `db:"id" json:"id"`
	LinkID     uint       `db:"link_id" json:"link_id"`
	ClickCount int64      `db:"click_count" json:"click_count"`
	Countries  Countries  `db:"countries" json:"countries"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type Countries []Country

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (c *Countries) Scan(value any) error {
	if value == nil {
		*c = []Country{}
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type: %T", value)
	}

	return json.Unmarshal(data, c)
}

func (c Countries) Value() (driver.Value, error) {
	if len(c) == 0 {
		return []byte("[]"), nil
	}

	return json.Marshal(c)
}

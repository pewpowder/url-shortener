package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

var LINK_SCHEMA = ` 
CREATE TABLE IF NOT EXISTS links (
	id SERIAL PRIMARY KEY,
	short_code VARCHAR(32) NOT NULL UNIQUE,
	original_url TEXT NOT NULL,
	expires_at TIMESTAMPTZ,
	is_active BOOLEAN NOT NULL DEFAULT true,
	is_private BOOLEAN NOT NULL DEFAULT false,
	password_hash VARCHAR(64),
	max_clicks INTEGER NOT NULL DEFAULT 0,
	tags JSONB NOT NULL DEFAULT '[]'::jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_links_short_code ON links(short_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_links_expires_at ON links(expires_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_links_deleted_at ON links(deleted_at);
CREATE INDEX IF NOT EXISTS idx_links_is_active ON links(is_active) WHERE deleted_at IS NULL;
`

type Link struct {
	ID           uint       `db:"id" json:"id"`
	ShortCode    string     `db:"short_code" json:"short_code"`
	OriginalURL  string     `db:"original_url" json:"original_url"`
	ExpiresAt    *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	IsPrivate    bool       `db:"is_private" json:"is_private"`
	PasswordHash *string    `db:"password_hash" json:"-"` // PasswordHash for private links. Can be null.
	MaxClicks    int        `db:"max_clicks" json:"max_clicks"`
	Tags         Tags       `db:"tags" json:"tags"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	// UserID       uint       `db:"user_id"`   // Integration with Auth-Service. (TODO: Add integration with auth service)
}

type Tags []string

func (t *Tags) Scan(value any) error {
	if value == nil {
		*t = []string{}
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type: %T", value)
	}

	return json.Unmarshal(data, t)
}

func (t Tags) Value() (driver.Value, error) {
	if len(t) == 0 {
		return []byte("[]"), nil
	}

	return json.Marshal(t)
}

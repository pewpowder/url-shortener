package entity

import "time"

var TAG_SCHEMA = `
CREATE TABLE IF NOT EXISTS tags (
	id SERIAL PRIMARY KEY,
	name VARCHAR(128) NOT NULL UNIQUE,
	color VARCHAR(7) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tags_deleted_at ON tags(deleted_at);
`

type Tag struct {
	ID        uint       `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Color     string     `db:"color" json:"color"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

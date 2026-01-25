package entity

import (
	"time"
)

var VISIT_SCHEMA = `
CREATE TABLE IF NOT EXISTS visits (
	id SERIAL PRIMARY KEY,
	link_id INTEGER NOT NULL,
	ip_address VARCHAR(45) NOT NULL,
	user_agent TEXT NOT NULL,
	referrer TEXT,
	country_code VARCHAR(2),
	timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,
	CONSTRAINT fk_visits_link FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_visits_link_id ON visits(link_id);
CREATE INDEX IF NOT EXISTS idx_visits_timestamp ON visits(timestamp);
CREATE INDEX IF NOT EXISTS idx_visits_country_code ON visits(country_code);
CREATE INDEX IF NOT EXISTS idx_visits_deleted_at ON visits(deleted_at);
CREATE INDEX IF NOT EXISTS idx_visits_link_timestamp ON visits(link_id, timestamp);
`

type Visit struct {
	ID          uint       `db:"id" json:"id"`
	LinkID      uint       `db:"link_id" json:"link_id"`
	IPAddress   string     `db:"ip_address" json:"ip_address"`
	UserAgent   string     `db:"user_agent" json:"user_agent"`
	Referrer    string     `db:"referrer" json:"referrer,omitempty"`
	CountryCode string     `db:"country_code" json:"country_code,omitempty"`
	Timestamp   time.Time  `db:"timestamp" json:"timestamp"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

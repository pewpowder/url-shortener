package container

import (
	"github.com/jmoiron/sqlx"
	"github.com/pewpowder/url-shortener/internal/config"
)

type Container interface {
	GetDB() *sqlx.DB
	GetConfig() *config.Config
	GetEnv() string
}

type container struct {
	db     *sqlx.DB
	config *config.Config
	env    string
}

func NewContainer(db *sqlx.DB, config *config.Config, env string) Container {
	return &container{
		db:     db,
		config: config,
		env:    env,
	}
}

func (c *container) GetDB() *sqlx.DB {
	return c.db
}

func (c *container) GetConfig() *config.Config {
	return c.config
}

func (c *container) GetEnv() string {
	return c.env
}

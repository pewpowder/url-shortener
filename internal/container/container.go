package container

import (
	"github.com/pewpowder/url-shortener/internal/config"
	"gorm.io/gorm"
)

type Container interface {
	GetDB() *gorm.DB
	GetConfig() *config.Config
	GetEnv() string
}

type container struct {
	db     *gorm.DB
	config *config.Config
	env    string
}

func NewContainer(repo *gorm.DB, config *config.Config, env string) Container {
	return &container{
		db:     repo,
		config: config,
		env:    env,
	}
}

func (c *container) GetDB() *gorm.DB {
	return c.db
}

func (c *container) GetConfig() *config.Config {
	return c.config
}

func (c *container) GetEnv() string {
	return c.env
}

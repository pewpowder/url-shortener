package container

import (
	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/repository"
)

type Container interface {
	GetRepository() repository.Repository
	GetConfig() *config.Config
	GetEnv() string
}

type container struct {
	repository repository.Repository
	config     *config.Config
	env        string
}

func NewContainer(repo repository.Repository, config *config.Config, env string) Container {
	return &container{
		repository: repo,
		config:     config,
		env:        env,
	}
}

func (c *container) GetRepository() repository.Repository {
	return c.repository
}

func (c *container) GetConfig() *config.Config {
	return c.config
}

func (c *container) GetEnv() string {
	return c.env
}

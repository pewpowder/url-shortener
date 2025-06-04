package service

import (
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
)

type LinkService interface{}

type linkService struct {
	container container.Container
}

func NewLinkService(container container.Container) LinkService {
	return &linkService{
		container,
	}
}

func (ls *linkService) CreateLink(data *dto.CreateLink) {
}

// TODO: Imeplement method for password determination with rules (like if IsPrivate is false I should reset password to nil etc)

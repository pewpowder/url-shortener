package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/service"
)

type LinkController interface {
	CreateLink(ctx *gin.Context)
	UpdateLink(ctx *gin.Context)
	DeleteLink(ctx *gin.Context)
}

type linkController struct {
	container container.Container
	service   service.LinkService
}

func NewLinkController(container container.Container) LinkController {
	return &linkController{
		container: container,
		service:   service.NewLinkService(container),
	}
}

func (lc *linkController) CreateLink(ctx *gin.Context) {
}

func (lc *linkController) UpdateLink(ctx *gin.Context) {
}

func (lc *linkController) DeleteLink(ctx *gin.Context) {
}

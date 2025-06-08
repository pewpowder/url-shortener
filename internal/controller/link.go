package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/service"
	"github.com/pewpowder/url-shortener/pkg/errors"
	"github.com/pewpowder/url-shortener/pkg/utils"
)

type LinkController interface {
	CreateLink(c *gin.Context)
	UpdateLink(c *gin.Context)
	DeleteLink(c *gin.Context)
	GetLinks(c *gin.Context)
	GetLinkDetails(c *gin.Context)
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

func (lc *linkController) GetLinks(c *gin.Context) {
	var query dto.LinkListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.InvalidQueryParams(c, err)
		return
	}
	query.Tags = utils.GinSplitString(query.Tags)

	links, err := lc.service.GetLinks(&query)

	if err != nil {
		return
	}

	linkList := dto.ToLinkList(links)

	c.JSON(200, linkList)
}

func (lc *linkController) GetLinkDetails(c *gin.Context) {
}

func (lc *linkController) CreateLink(c *gin.Context) {
	var body dto.CreateLink
	if err := c.ShouldBind(&body); err != nil {
		errors.InvalidBody(c, err)
		return
	}
}

func (lc *linkController) UpdateLink(c *gin.Context) {
}

func (lc *linkController) DeleteLink(c *gin.Context) {
}

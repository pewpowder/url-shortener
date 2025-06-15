package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/service"
	se "github.com/pewpowder/url-shortener/pkg/errors"
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
		se.InvalidQueryParams(c, err)
		return
	}
	if query.Tags != nil {
		tags := utils.GinSplitString(*query.Tags)
		query.Tags = &tags
	}

	links, err := lc.service.GetLinks(c.Request.Context(), query)
	if err != nil {
		se.HandleError(c, err)
		return
	}

	linkList := dto.ToLinkList(links)

	c.JSON(200, linkList)
}

func (lc *linkController) GetLinkDetails(c *gin.Context) {
}

func (lc *linkController) CreateLink(c *gin.Context) {
	var data dto.CreateLink
	if err := c.ShouldBind(&data); err != nil {
		se.InvalidRequestData(c, err)
		return
	}

	link, err := lc.service.CreateLink(c.Request.Context(), data)
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToLinkDetails(link))
}

func (lc *linkController) UpdateLink(c *gin.Context) {
}

func (lc *linkController) DeleteLink(c *gin.Context) {
}

// TODO: Determine where i should set user_id

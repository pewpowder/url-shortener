package controller

import (
	"net/http"
	"strconv"

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
	PatchLink(c *gin.Context)
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
		service:   service.NewLinkService(container, service.NewTagService(container)),
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

	c.JSON(200, dto.ToLinkList(links)) // TODO: return by default structured object (data, total, page, size etc) instead of null
}

func (lc *linkController) GetLinkDetails(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		se.InvalidRequestData(c, err)
		return
	}

	link, err := lc.service.GetLinkById(c.Request.Context(), uint(id64))
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToLinkDetails(link))
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
	var data dto.CreateLink
	if err := c.ShouldBind(&data); err != nil {
		se.InvalidRequestData(c, err)
		return
	}

	link, err := lc.service.UpdateLink(c.Request.Context(), data)
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToLinkDetails(link))
}

func (lc *linkController) PatchLink(c *gin.Context) {
	var data dto.PatchLink
	if err := c.ShouldBind(&data); err != nil {
		se.InvalidRequestData(c, err)
		return
	}

	link, err := lc.service.PatchLink(c.Request.Context(), data)
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToLinkDetails(link))
}

func (lc *linkController) DeleteLink(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		se.InvalidRequestData(c, err)
		return
	}

	err = lc.service.DeleteLink(c.Request.Context(), uint(id64))
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"success": true, "id": id64, "message": "Link deleted successfully"})
}

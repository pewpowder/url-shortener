package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/service"
	se "github.com/pewpowder/url-shortener/pkg/errors"
)

type TagController interface {
	CreateTag(c *gin.Context)
	UpdateTag(c *gin.Context)
	DeleteTag(c *gin.Context)
	GetTags(c *gin.Context)
	GetTagByID(c *gin.Context)
}

type tagController struct {
	container container.Container
	service   service.TagService
}

func NewTagController(container container.Container) TagController {
	return &tagController{
		container: container,
		service:   service.NewTagService(container),
	}
}

func (tc *tagController) CreateTag(c *gin.Context) {
	var tagDto dto.TagRequest
	if err := c.ShouldBind(&tagDto); err != nil {
		se.InvalidRequestData(c, err)
		return
	}

	tag, err := tc.service.CreateTag(tagDto)
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToTagResponse(tag))
}

func (tc *tagController) UpdateTag(c *gin.Context) {
}

func (tc *tagController) DeleteTag(c *gin.Context) {
}

func (tc *tagController) GetTags(c *gin.Context) {
	tags, err := tc.service.GetTags()
	if err != nil {
		se.HandleError(c, err)
		return
	}

	tagDtos := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		tagDtos = append(tagDtos, dto.ToTagResponse(tag))
	}

	c.JSON(http.StatusOK, tagDtos)
}

func (tc *tagController) GetTagByID(c *gin.Context) {
	var id uint
	if err := c.ShouldBindUri(&id); err != nil {
		se.HandleError(c, err)
		return
	}

	tag, err := tc.service.GetTagByID(id)
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTagResponse(tag))
}

package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/service"
	se "github.com/pewpowder/url-shortener/pkg/errors"
)

type TagController interface {
	CreateTag(c *gin.Context)
	UpdateTag(c *gin.Context)
	PatchTag(c *gin.Context)
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
	var tagDto dto.CreateOrUpdateTag
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
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 {
		se.InvalidQueryParams(c, err)
		return
	}

	var tagDto dto.CreateOrUpdateTag
	if err := c.ShouldBind(&tagDto); err != nil {
		se.InvalidRequestData(c, err)
		return
	}

	tag, err := tc.service.UpdateTag(uint(id64), tagDto)
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTagResponse(tag))
}

func (tc *tagController) PatchTag(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 {
		se.InvalidQueryParams(c, err)
		return
	}

	var tagDto dto.PatchTag
	if err := c.ShouldBind(&tagDto); err != nil {
		se.InvalidRequestData(c, err)
		return
	}

	tag, err := tc.service.PatchTag(uint(id64), tagDto)
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTagResponse(tag))
}

func (tc *tagController) DeleteTag(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 {
		se.InvalidQueryParams(c, err)
		return
	}

	tag, err := tc.service.DeleteTag(uint(id64))
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTagResponse(tag))
}

func (tc *tagController) GetTags(c *gin.Context) {
	var q dto.TagListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		se.InvalidQueryParams(c, err)
		return
	}

	tags, err := tc.service.GetTags(q)
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
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 {
		se.InvalidQueryParams(c, err)
		return
	}

	tag, err := tc.service.GetTagByID(uint(id64))
	if err != nil {
		se.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTagResponse(tag))
}

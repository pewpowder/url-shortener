package dto

import (
	"time"

	"github.com/pewpowder/url-shortener/internal/entity"
)

type TagFilter struct {
	Name  *string  `form:"filter[name]"`
	Names []string `form:"filter[names]"`
	Color *string  `form:"filter[color]"`
}

type TagListQuery struct {
	Pagination Pagination
	Sort       Sort // TODO: make optional
	Filter     *TagFilter
}

type CreateOrUpdateTag struct {
	Name  string `json:"name" binding:"required,alphanum,max=50"`
	Color string `json:"color" binding:"omitempty,hexcolor"`
}

type PatchTag struct {
	Name  *string `json:"name" binding:"omitempty,alphanum,max=50"`
	Color *string `json:"color" binding:"omitempty,hexcolor"`
}

type TagResponse struct {
	ID        uint       `json:"id" validate:"required"`
	CreatedAt time.Time  `json:"created_at" validate:"required"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" validate:"required,alphanum,max=50"`
	Color     string     `json:"color" validate:"required,hexcolor"`
}

type TagEmbedding struct {
	ID    uint   `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,alphanum,max=50"`
	Color string `json:"color" validate:"required,hexcolor"`
}

func ToTagResponse(tag entity.Tag) TagResponse {
	return TagResponse{
		ID:        tag.ID,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
		Name:      tag.Name,
		Color:     tag.Color,
	}
}

func ToTagEmbedding(tag entity.Tag) TagEmbedding {
	return TagEmbedding{
		ID:    tag.ID,
		Name:  tag.Name,
		Color: tag.Color,
	}
}

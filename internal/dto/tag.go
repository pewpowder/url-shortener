package dto

import (
	"time"

	"github.com/pewpowder/url-shortener/internal/entity"
)

type TagRequest struct {
	Name  string `json:"name" binding:"required,alphanum,max=50"`
	Color string `json:"color" binding:"required,hexcolor"`
}

type TagResponse struct {
	ID        uint      `json:"id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	Name      string    `json:"name" validate:"required,alphanum,max=50"`
	Color     string    `json:"color" validate:"required,hexcolor"`
}

type TagEmbedding struct {
	ID    uint   `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,alphanum,max=50"`
	Color string `json:"color" validate:"required,hexcolor"`
}

func FromTagRequest(tag TagRequest) entity.Tag {
	return entity.Tag{
		Name:  tag.Name,
		Color: tag.Color,
	}
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

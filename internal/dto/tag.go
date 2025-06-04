package dto

import "github.com/pewpowder/url-shortener/internal/entity"

type Tag struct {
	Name  string `json:"name" binding:"required,alphanum,max=50" validate:"required,alphanum,max=50"`
	Color string `json:"color" binding:"required,hexcolor" validate:"required,hexcolor"`
}

func ToTag(tag entity.Tag) Tag {
	return Tag{
		Name:  tag.Name,
		Color: tag.Color,
	}
}

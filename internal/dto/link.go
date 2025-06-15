package dto

import (
	"errors"
	"time"

	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/pkg/utils"
)

type CreateLink struct {
	URL       string   `json:"url" binding:"required,url"`
	IsPrivate bool     `json:"is_private" binding:"default=false"`
	IsActive  bool     `json:"is_active" binding:"default=true"`
	Password  string   `json:"password" binding:"required_if=is_private true min=6,containsany=!@#$%^&*"`
	MaxClicks int      `json:"max_clicks" binding:"gte=0,default=0"` // 0 means no limit
	Tags      []string `json:"tags" binding:"max=10"`
	ExpiresAt *string  `json:"expires_at" binding:"datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateLink struct {
	URL       *string   `json:"url" binding:"url"`
	IsPrivate *bool     `json:"is_private"`
	IsActive  *bool     `json:"is_active"`
	Password  *string   `json:"password" binding:"required_if=is_private true min=6,containsany=!@#$%^&*"`
	MaxClicks *int      `json:"max_clicks" binding:"gte=0"`
	Tags      *[]string `json:"tags" binding:"max=10"`
	ExpiresAt *string   `json:"expires_at" binding:"datetime=2006-01-02T15:04:05Z07:00"`
}

type LinkListQuery struct {
	Tags      *[]string `form:"tags"`
	IsPrivate *bool     `form:"is_private"`
	IsActive  *bool     `form:"is_active"`
	Expired   *bool     `form:"expired"`
}

// DTO for response validates only in development mode
type LinkList struct {
	ID          uint           `json:"id" validate:"required"`
	ShortCode   string         `json:"short_code" validate:"required"`
	OriginalURL string         `json:"original_url" validate:"required"`
	IsActive    bool           `json:"is_active" validate:"required"`
	Tags        []TagEmbedding `json:"tags"`
	CreatedAt   string         `json:"created_at" validate:"required"`
	UpdatedAt   string         `json:"updated_at" validate:"required"`
	ExpiresAt   *string        `json:"expires_at"`
}

type LinkDetails struct {
	ID          uint           `json:"id" validate:"required"`
	ShortCode   string         `json:"short_code" validate:"required"`
	OriginalURL string         `json:"original_url" validate:"required"`
	IsActive    bool           `json:"is_active" validate:"required"`
	IsPrivate   bool           `json:"is_private" validate:"required"`
	UserID      uint           `json:"user_id"` // validate:"required"
	MaxClicks   int            `json:"max_clicks" validate:"required"`
	Tags        []TagEmbedding `json:"tags"`
	CreatedAt   string         `json:"created_at" validate:"required"`
	UpdatedAt   string         `json:"updated_at" validate:"required"`
	ExpiresAt   *string        `json:"expires_at"`
	DeletedAt   *string        `json:"deleted_at"`
}

func FromCreateLink(createLink CreateLink, passHash *string, tags []entity.Tag) (entity.Link, error) {
	expiresAt, _ := utils.ParseOptionalTime(createLink.ExpiresAt)

	if expiresAt != nil && expiresAt.Before(time.Now()) {
		return entity.Link{}, errors.New("expires_at must be non past value")
	}

	return entity.Link{
		OriginalURL:  createLink.URL,
		IsPrivate:    createLink.IsPrivate,
		IsActive:     createLink.IsActive,
		PasswordHash: passHash,
		MaxClicks:    createLink.MaxClicks,
		Tags:         tags,
		ExpiresAt:    expiresAt,
	}, nil
}

func FromUpdateLink(updateLink UpdateLink, passHash *string, tags []entity.Tag) (entity.Link, error) {
	link := entity.Link{}

	if updateLink.URL != nil {
		link.OriginalURL = *updateLink.URL
	}

	if updateLink.IsPrivate != nil {
		link.IsPrivate = *updateLink.IsPrivate
	}

	if updateLink.IsActive != nil {
		link.IsActive = *updateLink.IsActive
	}

	if updateLink.MaxClicks != nil {
		link.MaxClicks = *updateLink.MaxClicks
	}

	if tags != nil {
		link.Tags = tags
	}

	if updateLink.ExpiresAt != nil {
		parsedTime, _ := utils.ParseOptionalTime(updateLink.ExpiresAt)
		if parsedTime != nil && parsedTime.Before(time.Now()) {
			return entity.Link{}, errors.New("expires_at must be non past value")
		}

		link.ExpiresAt = parsedTime
	}

	link.PasswordHash = passHash

	return link, nil
}

func ToLinkList(links []entity.Link) []LinkList {
	var linkList []LinkList

	for _, link := range links {
		tags := make([]TagEmbedding, len(link.Tags))

		for _, tag := range link.Tags {
			tags = append(tags, ToTagEmbedding(tag))
		}

		expiresAt := link.ExpiresAt.Format(config.TIME_FORMAT)
		createdAt := link.CreatedAt.Format(config.TIME_FORMAT)

		item := LinkList{
			ShortCode:   link.ShortCode,
			OriginalURL: link.OriginalURL,
			IsActive:    link.IsActive,
			CreatedAt:   createdAt,
			Tags:        tags,
			ExpiresAt:   &expiresAt,
		}

		linkList = append(linkList, item)
	}

	return linkList
}

func ToLinkDetails(link entity.Link) LinkDetails {
	tags := make([]TagEmbedding, len(link.Tags))
	for _, tag := range link.Tags {
		tags = append(tags, ToTagEmbedding(tag))
	}

	var deletedAt string
	if link.DeletedAt.Valid {
		deletedAt = link.DeletedAt.Time.Format(config.TIME_FORMAT)
	}

	var expiresAt *string
	if link.ExpiresAt != nil {
		formatted := link.ExpiresAt.Format(config.TIME_FORMAT)
		expiresAt = &formatted
	}

	createdAt := link.CreatedAt.Format(config.TIME_FORMAT)
	updatedAt := link.UpdatedAt.Format(config.TIME_FORMAT)

	return LinkDetails{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		IsActive:    link.IsActive,
		IsPrivate:   link.IsPrivate,
		UserID:      link.UserID,
		MaxClicks:   link.MaxClicks,
		Tags:        tags,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		ExpiresAt:   expiresAt,
		DeletedAt:   &deletedAt,
	}
}

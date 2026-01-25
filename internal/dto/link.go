package dto

import (
	"time"

	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/pkg/utils"
)

type CreateLink struct {
	URL       string   `json:"url" binding:"required,url"`
	IsPrivate *bool    `json:"is_private"`
	IsActive  *bool    `json:"is_active"`
	Password  string   `json:"password" binding:"omitempty,required_if=IsPrivate true,min=6,containsany=!@#$%^&*"`
	MaxClicks int      `json:"max_clicks" binding:"gte=0,lte=100000"` // 0 means no limit
	Tags      []string `json:"tags" binding:"max=10"`
	ExpiresAt *string  `json:"expires_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type PatchLink struct {
	URL       *string  `json:"url" binding:"omitempty,url"`
	IsPrivate *bool    `json:"is_private"`
	IsActive  *bool    `json:"is_active"`
	Password  *string  `json:"password" binding:"omitempty,required_if=IsPrivate true,min=6,containsany=!@#$%^&*"`
	MaxClicks *int     `json:"max_clicks" binding:"omitempty,gte=0,lte=100000"`
	Tags      []string `json:"tags" binding:"omitempty,max=10"`
	ExpiresAt *string  `json:"expires_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type LinkListQuery struct {
	Tags      *[]string `form:"tags"`
	IsPrivate *bool     `form:"is_private"`
	IsActive  *bool     `form:"is_active"`
	Expired   *bool     `form:"expired"`
}

// DTO for response validates only in development mode
type LinkList struct {
	ID          uint     `json:"id" validate:"required"`
	ShortCode   string   `json:"short_code" validate:"required"`
	OriginalURL string   `json:"original_url" validate:"required"`
	IsActive    bool     `json:"is_active" validate:"required"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at" validate:"required"`
	UpdatedAt   string   `json:"updated_at" validate:"required"`
	ExpiresAt   *string  `json:"expires_at"`
}

type LinkDetails struct {
	ID          uint     `json:"id" validate:"required"`
	ShortCode   string   `json:"short_code" validate:"required"`
	OriginalURL string   `json:"original_url" validate:"required"`
	IsActive    bool     `json:"is_active" validate:"required"`
	IsPrivate   bool     `json:"is_private" validate:"required"`
	MaxClicks   int      `json:"max_clicks" validate:"required"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at" validate:"required"`
	UpdatedAt   *string  `json:"updated_at" validate:"required"`
	ExpiresAt   *string  `json:"expires_at"`
	DeletedAt   *string  `json:"deleted_at"`
}

func FromCreateLink(createLink CreateLink, shortCode string, passHash *string, tags []string, expiresAt *time.Time) (entity.Link, error) {
	return entity.Link{
		ShortCode:    shortCode,
		OriginalURL:  createLink.URL,
		IsPrivate:    utils.DerefBool(createLink.IsPrivate, false),
		IsActive:     createLink.IsActive != nil && *createLink.IsActive,
		PasswordHash: passHash,
		MaxClicks:    createLink.MaxClicks,
		Tags:         tags,
		ExpiresAt:    expiresAt,
	}, nil
}

// 1. I should get link from db
// 2. I should update link with new data
// 3. I should update link in db
// 4. I should return link
// 5. Is this a good approach?
func FromUpdateLink(updateLink PatchLink, passHash *string, tags []string, expiresAt *time.Time) (entity.Link, error) {
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

	if expiresAt != nil {
		link.ExpiresAt = expiresAt
	}

	link.PasswordHash = passHash

	return link, nil
}

func ToLinkList(links []entity.Link) []LinkList {
	var linkList []LinkList

	for _, link := range links {
		var expiresAt *string
		if link.ExpiresAt != nil {
			v := link.ExpiresAt.Format(config.TIME_FORMAT)
			expiresAt = &v
		}

		item := LinkList{
			ID:          link.ID,
			ShortCode:   link.ShortCode,
			OriginalURL: link.OriginalURL,
			IsActive:    link.IsActive,
			Tags:        link.Tags,
			CreatedAt:   link.CreatedAt.Format(config.TIME_FORMAT),
			UpdatedAt:   link.UpdatedAt.Format(config.TIME_FORMAT),
			ExpiresAt:   expiresAt,
		}

		linkList = append(linkList, item)
	}

	return linkList
}

func ToLinkDetails(link entity.Link) LinkDetails {
	var deletedAt *string
	if link.DeletedAt != nil {
		formatted := link.DeletedAt.Format(config.TIME_FORMAT)
		deletedAt = &formatted
	}

	var expiresAt *string
	if link.ExpiresAt != nil {
		formatted := link.ExpiresAt.Format(config.TIME_FORMAT)
		expiresAt = &formatted
	}

	var updatedAt *string
	if link.UpdatedAt != nil {
		formatted := link.UpdatedAt.Format(config.TIME_FORMAT)
		updatedAt = &formatted
	}

	return LinkDetails{
		ID:          link.ID,
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		IsActive:    link.IsActive,
		IsPrivate:   link.IsPrivate,
		MaxClicks:   link.MaxClicks,
		Tags:        link.Tags,
		CreatedAt:   link.CreatedAt.Format(config.TIME_FORMAT),
		UpdatedAt:   updatedAt,
		ExpiresAt:   expiresAt,
		DeletedAt:   deletedAt,
	}
}

package dto

import (
	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/pkg/utils"
)

type CreateLink struct {
	URL       string   `json:"url" binding:"required,url"`
	IsPrivate *bool    `json:"is_private"`
	IsActive  *bool    `json:"is_active"`                                                                 // Optional, defaults to true
	Password  string   `json:"password" binding:"required_if=is_private true min=6,containsany=!@#$%^&*"` // Optional for private links
	MaxClicks int      `json:"max_clicks" binding:"gte=0"`                                                // Optional, 0 means no limit
	Tags      []string `json:"tags" binding:"max=10"`                                                     // Optional, tags for the link
	ExpiresAt *string  `json:"expires_at" binding:"datetime=2006-01-02T15:04:05Z07:00"`                   // Optional, ISO 8601 format
}

type UpdateLink struct {
	URL       string    `json:"url" binding:"url"`
	IsPrivate *bool     `json:"is_private"`
	IsActive  *bool     `json:"is_active"`
	Password  string    `json:"password" binding:"required_if=is_private true min=6,containsany=!@#$%^&*"`
	MaxClicks *int      `json:"max_clicks" binding:"gte=0"`
	Tags      *[]string `json:"tags" binding:"max=10"`
	ExpiresAt *string   `json:"expires_at" binding:"datetime=2006-01-02T15:04:05Z07:00"`
}

// DTO for retrive validates only in development mode
type LinkList struct {
	ShortCode   string  `json:"short_code" validate:"required"`
	OriginalURL string  `json:"original_url" validate:"required"`
	IsActive    bool    `json:"is_active" validate:"required"`
	Tags        []Tag   `json:"tags"`
	ExpiresAt   *string `json:"expires_at"`
	CreatedAt   *string `json:"created_at" validate:"required"`
}

type LinkDetails struct {
	ShortCode   string  `json:"short_code" validate:"required"`
	OriginalURL string  `json:"original_url" validate:"required"`
	IsActive    bool    `json:"is_active" validate:"required"`
	IsPrivate   bool    `json:"is_private" validate:"required"`
	UserID      uint    `json:"user_id"` // validate:"required"
	MaxClicks   int     `json:"max_clicks" validate:"required"`
	Tags        []Tag   `json:"tags"`
	ExpiresAt   *string `json:"expires_at"`
	CreatedAt   *string `json:"created_at" validate:"required"`
	UpdatedAt   *string `json:"updated_at" validate:"required"`
	DeletedAt   *string `json:"deleted_at"`
}

func FromCreateLink(createLink *CreateLink, passHash *string, tags []entity.Tag) (*entity.Link, error) {
	expiresAt, _ := utils.ParseOptionalTime(createLink.ExpiresAt) // validates with DTO tags

	return &entity.Link{
		OriginalURL:  createLink.URL,
		IsPrivate:    utils.DerefBool(createLink.IsPrivate, false),
		IsActive:     utils.DerefBool(createLink.IsActive, true),
		PasswordHash: passHash,
		MaxClicks:    createLink.MaxClicks,
		Tags:         tags,
		ExpiresAt:    expiresAt, // TODO: check that time is not the past
	}, nil
}

func FromUpdateLink(updateLink *UpdateLink, passHash *string, tags []entity.Tag) *entity.Link {
	link := &entity.Link{}

	if updateLink.URL != "" {
		link.OriginalURL = updateLink.URL
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
		parsedTime, _ := utils.ParseOptionalTime(updateLink.ExpiresAt) // validates with DTO tags
		link.ExpiresAt = parsedTime
	}

	link.PasswordHash = passHash

	return link
}

func ToLinkList(link *entity.Link) *LinkList {
	tags := make([]Tag, len(link.Tags))

	for _, tag := range link.Tags {
		tags = append(tags, ToTag(tag))
	}

	expiresAt := link.ExpiresAt.Format(config.TIME_FORMAT)
	createdAt := link.CreatedAt.Format(config.TIME_FORMAT)

	return &LinkList{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		ExpiresAt:   &expiresAt,
		IsActive:    link.IsActive,
		CreatedAt:   &createdAt,
		Tags:        tags,
	}
}

func ToLinkDetails(link *entity.Link) *LinkDetails {
	tags := make([]Tag, len(link.Tags))
	for _, tag := range link.Tags {
		tags = append(tags, ToTag(tag))
	}

	var deletedAt string
	if link.DeletedAt.Valid {
		deletedAt = link.DeletedAt.Time.Format(config.TIME_FORMAT)
	}

	expiresAt := link.ExpiresAt.Format(config.TIME_FORMAT)
	createdAt := link.CreatedAt.Format(config.TIME_FORMAT)
	updatedAt := link.UpdatedAt.Format(config.TIME_FORMAT)

	return &LinkDetails{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		IsActive:    link.IsActive,
		IsPrivate:   link.IsPrivate,
		UserID:      link.UserID,
		MaxClicks:   link.MaxClicks,
		Tags:        tags,
		ExpiresAt:   &expiresAt,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
		DeletedAt:   &deletedAt,
	}
}

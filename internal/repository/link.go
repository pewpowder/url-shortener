package repository

import (
	"fmt"

	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	se "github.com/pewpowder/url-shortener/pkg/errors"
	"gorm.io/gorm"
)

type LinkRepository interface {
	GetLinks(query dto.LinkListQuery) ([]entity.Link, error)
	GetLinkById(id uint) (entity.Link, error)
	GetLinkByShortCode(code string) (entity.Link, error)
}

type linkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
	return &linkRepository{
		db: db,
	}
}

func (lr *linkRepository) GetLinks(query dto.LinkListQuery) ([]entity.Link, error) {
	var links []entity.Link
	db := lr.db.Model(&entity.Link{})

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.IsPrivate != nil {
		db = db.Where("is_private = ?", *query.IsPrivate)
	}

	// if query.UserID != nil {
	// 	db = db.Where("user_id = ?", *query.UserID)
	// }

	if query.Expired != nil {
		if *query.Expired {
			db = db.Where("expires_at < NOW()")
		} else {
			db = db.Where("expires_at > NOW()")
		}
	}

	if query.Tags != nil && len(*query.Tags) > 0 {
		db = db.
			Joins("JOIN link_tags ON link_tags.link_id = links.id").
			Joins("JOIN tags ON tags.id = link_tags.tag_id").
			Where("tags.name IN ?", *query.Tags)
	}

	if err := db.Preload("Tags").Find(&links).Error; err != nil {
		return nil, se.NewServiceError(fmt.Sprintf("failed to get links: %s", err), se.GormErrorToErrorType(err), err)
	}

	return links, nil
}

func (lr *linkRepository) GetLinkById(id uint) (entity.Link, error) {
	return lr.getLinkBy("id = ?", id)
}

func (lr *linkRepository) GetLinkByShortCode(code string) (entity.Link, error) {
	return lr.getLinkBy("short_code = ?", code)
}

func (lr *linkRepository) getLinkBy(where string, args ...any) (entity.Link, error) {
	var link entity.Link

	if err := lr.db.Where(where, args...).Preload("Tags").Find(&link).Error; err != nil {
		return entity.Link{}, se.NewServiceError(fmt.Sprintf("failed to get link: %s", err), se.GormErrorToErrorType(err), err)
	}

	return link, nil
}

package repository

import (
	"fmt"

	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	se "github.com/pewpowder/url-shortener/pkg/errors"
	"gorm.io/gorm"
)

// TODO: change link accordingly to the tag (also remove verbose methods and pass entity.Link to repo instead of dto.Link)
type LinkRepository interface {
	GetLinks(query dto.LinkListQuery) ([]entity.Link, error)
	GetLinkById(id uint) (entity.Link, error)
	GetLinkByShortCode(code string) (entity.Link, error)
	CreateLink(link entity.Link) (entity.Link, error)
	UpdateLink(link entity.Link) (entity.Link, error)
	DeleteLink(id uint) error
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

	// if query.UserID != nil {
	// 	db = db.Where("user_id = ?", *query.UserID)
	// }

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.IsPrivate != nil {
		db = db.Where("is_private = ?", *query.IsPrivate)
	}

	if query.Expired != nil {
		if *query.Expired {
			db = db.Where("expires_at < NOW()")
		} else {
			db = db.Where("expires_at > NOW() OR expires_at IS NULL")
		}
	}

	if query.Tags != nil && len(*query.Tags) > 0 {
		db = db.
			Joins("JOIN link_tags ON link_tags.link_id = links.id").
			Joins("JOIN tags ON tags.id = link_tags.tag_id").
			Where("tags.name IN ?", *query.Tags).
			Distinct()
	}

	if err := db.Preload("Tags").Find(&links).Error; err != nil {
		return nil, se.NewServiceErrorFromGorm("failed to get links", err)
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

	if err := lr.db.Where(where, args...).Preload("Tags").First(&link).Error; err != nil {
		return entity.Link{}, se.NewServiceErrorFromGorm("failed to get link", err)
	}

	return link, nil
}

func (lr *linkRepository) CreateLink(link entity.Link) (entity.Link, error) {
	db := lr.db.Model(&entity.Link{})

	if err := db.Create(&link).Error; err != nil {
		return entity.Link{}, se.NewServiceErrorFromGorm("failed to create link", err)
	}

	return link, nil
}

func (lr *linkRepository) UpdateLink(link entity.Link) (entity.Link, error) {
	db := lr.db.Model(&entity.Link{})

	if err := db.Updates(&link).Error; err != nil {
		return entity.Link{}, se.NewServiceErrorFromGorm(fmt.Sprintf("failed to update link with id = %d", link.ID), err)
	}

	return link, nil
}

func (lr *linkRepository) DeleteLink(id uint) error {
	db := lr.db.Model(&entity.Link{})

	if err := db.Delete(&entity.Link{}, id).Error; err != nil {
		return se.NewServiceErrorFromGorm(fmt.Sprintf("failed to delete link with id = %d", id), err)
	}

	return nil
}

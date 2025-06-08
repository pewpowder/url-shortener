package repository

import (
	"github.com/pewpowder/url-shortener/internal/entity"
	"gorm.io/gorm"
)

const (
	FIND_LINKS_QUERY = "IsPrivate = ? AND "
)

type LinkRepository interface {
	GetLinks(query string, where ...any) ([]entity.Link, error)
}

type linkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
	return &linkRepository{
		db: db,
	}
}

func (lr *linkRepository) GetLinks(query string, where ...any) ([]entity.Link, error) {
	var links []entity.Link

	tx := lr.db.Where(query, where...).Find(&links)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return links, nil
}

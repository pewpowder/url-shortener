package repository

import (
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	"gorm.io/gorm"
)

type TagRepository interface {
	CreateTag(createTag dto.TagRequest) (entity.Tag, error)
	GetTags() ([]entity.Tag, error)
	GetTagByID(id uint) (entity.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{
		db: db,
	}
}

func (tr *tagRepository) GetTable() string {
	return "tags"
}

func (tr *tagRepository) CreateTag(tagDto dto.TagRequest) (entity.Tag, error) {
	tag := dto.FromTagRequest(tagDto)

	if err := tr.db.Create(tag).Error; err != nil {
		return entity.Tag{}, err
	}

	return tag, nil
}

func (tr *tagRepository) GetTags() ([]entity.Tag, error) {
	var tags []entity.Tag

	if err := tr.db.Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (tr *tagRepository) GetTagByID(id uint) (entity.Tag, error) {
	var tag entity.Tag

	if err := tr.db.First(&tag, id).Error; err != nil {
		return entity.Tag{}, err
	}

	return tag, nil
}

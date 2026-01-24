package repository

import (
	"fmt"

	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	se "github.com/pewpowder/url-shortener/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagRepository interface {
	CreateTag(createTag entity.Tag) (entity.Tag, error)
	UpdateTag(id uint, updateTag entity.Tag) (entity.Tag, error)
	PatchTag(id uint, patchTag entity.Tag) (entity.Tag, error)
	DeleteTag(id uint) (entity.Tag, error)
	GetTags(q dto.TagListQuery) ([]entity.Tag, error)
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

func (tr *tagRepository) CreateTag(tag entity.Tag) (entity.Tag, error) {
	if err := tr.db.Create(&tag).Error; err != nil {
		return entity.Tag{}, se.NewServiceErrorFromGorm("failed to create tag", err)
	}

	return tag, nil
}

func (tr *tagRepository) UpdateTag(id uint, tag entity.Tag) (entity.Tag, error) {
	if err := tr.db.Save(&tag).Error; err != nil {
		return entity.Tag{}, se.NewServiceErrorFromGorm(fmt.Sprintf("failed to update tag with id %d", id), err)
	}

	return tag, nil
}

func (tr *tagRepository) PatchTag(id uint, tag entity.Tag) (entity.Tag, error) {
	out := entity.Tag{}

	tx := tr.db.
		Model(&entity.Tag{}).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(tag).
		Scan(&out)

	if tx.Error != nil {
		return entity.Tag{}, se.NewServiceErrorFromGorm(fmt.Sprintf("failed to update tag with id %d", id), tx.Error)
	}

	if tx.RowsAffected == 0 {
		return entity.Tag{}, se.NewServiceErrorFromGorm(fmt.Sprintf("tag with id %d not found", id), gorm.ErrRecordNotFound)
	}

	return out, nil
}

func (tr *tagRepository) DeleteTag(id uint) (entity.Tag, error) {
	out := entity.Tag{}
	tx := tr.db.Clauses(clause.Returning{}).Delete(&out, id)

	if tx.Error != nil {
		return entity.Tag{}, se.NewServiceErrorFromGorm(fmt.Sprintf("failed to delete tag with id %d", id), tx.Error)
	}

	if tx.RowsAffected == 0 {
		return entity.Tag{}, se.NewServiceErrorFromGorm(fmt.Sprintf("tag with id %d not found", id), gorm.ErrRecordNotFound)
	}

	return out, nil
}

func (tr *tagRepository) GetTags(q dto.TagListQuery) ([]entity.Tag, error) {
	sqlQuery := tr.db.Model(&entity.Tag{})
	sqlQuery = sqlQuery.Offset(q.Pagination.GetOffset()).Limit(q.Pagination.Size)
	if isSortableField(q.Sort.By) {
		orderClause := clause.OrderByColumn{Column: clause.Column{Name: q.Sort.By}, Desc: q.Sort.Order == "desc"}
		sqlQuery = sqlQuery.Order(orderClause)
	}

	if q.Filter != nil {
		if len(q.Filter.Names) > 0 {
			sqlQuery = sqlQuery.Where("name IN ?", q.Filter.Names)
		} else if q.Filter.Name != nil {
			sqlQuery = sqlQuery.Where("name = ?", *q.Filter.Name)
		}

		if q.Filter.Color != nil {
			sqlQuery = sqlQuery.Where("color = ?", *q.Filter.Color)
		}
	}

	var tags []entity.Tag
	if err := sqlQuery.Find(&tags).Error; err != nil {
		return nil, se.NewServiceErrorFromGorm("failed to get tags", err)
	}

	return tags, nil
}

func (tr *tagRepository) GetTagByID(id uint) (entity.Tag, error) {
	var tag entity.Tag

	if err := tr.db.First(&tag, id).Error; err != nil {
		return entity.Tag{}, se.NewServiceErrorFromGorm(fmt.Sprintf("failed to get tag by id: %d", id), err)
	}

	return tag, nil
}

func isSortableField(field string) bool {
	sortableFields := map[string]bool{
		"name":       true,
		"color":      true,
		"created_at": true,
		"updated_at": true,
	}

	return sortableFields[field]
}

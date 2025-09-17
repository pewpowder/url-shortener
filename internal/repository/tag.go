package repository

import (
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	se "github.com/pewpowder/url-shortener/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagRepository interface {
	// TODO: change link accordingly to the tag (also remove verbose methods and pass entity.Link to repo instead of dto.Link)
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
		return entity.Tag{}, se.NewServiceError("failed to create tag", se.GormErrorToErrorType(err), err)
	}

	return tag, nil
}

func (tr *tagRepository) UpdateTag(id uint, tag entity.Tag) (entity.Tag, error) {
	if err := tr.db.Save(&tag).Error; err != nil {
		return entity.Tag{}, se.NewServiceError("failed to update tag", se.GormErrorToErrorType(err), err)
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
		return entity.Tag{}, se.NewServiceError("failed to update tag", se.GormErrorToErrorType(tx.Error), tx.Error)
	}

	if tx.RowsAffected == 0 {
		return entity.Tag{}, se.NewServiceError("tag not found", se.ErrTypeNotFound, nil)
	}

	return out, nil
}

func (tr *tagRepository) DeleteTag(id uint) (entity.Tag, error) {
	out := entity.Tag{}
	tx := tr.db.Clauses(clause.Returning{}).Delete(&out, id)

	if tx.Error != nil {
		return entity.Tag{}, se.NewServiceError("failed to delete tag", se.GormErrorToErrorType(tx.Error), tx.Error)
	}

	if tx.RowsAffected == 0 {
		return entity.Tag{}, se.NewServiceError("tag not found", se.ErrTypeNotFound, nil)
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
		if q.Filter.Name != nil {
			sqlQuery = sqlQuery.Where("name = ?", *q.Filter.Name)
		}

		if q.Filter.Color != nil {
			sqlQuery = sqlQuery.Where("color = ?", *q.Filter.Color)
		}
	}

	var tags []entity.Tag
	if err := sqlQuery.Find(&tags).Error; err != nil {
		return nil, se.NewServiceError("failed to get tags", se.GormErrorToErrorType(err), err)
	}

	return tags, nil
}

func (tr *tagRepository) GetTagByID(id uint) (entity.Tag, error) {
	var tag entity.Tag

	if err := tr.db.First(&tag, id).Error; err != nil {
		return entity.Tag{}, se.NewServiceError("failed to get tag by id", se.GormErrorToErrorType(err), err)
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

package service

import (
	"time"

	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/internal/repository"
)

const DEFAULT_COLOR = "#4f46e5"

type TagService interface {
	CreateTag(createTag dto.CreateOrUpdateTag) (entity.Tag, error)
	UpdateTag(id uint, updateTag dto.CreateOrUpdateTag) (entity.Tag, error)
	PatchTag(id uint, patchTag dto.PatchTag) (entity.Tag, error)
	DeleteTag(id uint) (entity.Tag, error)
	GetTags(q dto.TagListQuery) ([]entity.Tag, error)
	GetTagByID(id uint) (entity.Tag, error)
}

type tagService struct {
	container container.Container
	repo      repository.TagRepository
}

func NewTagService(container container.Container) TagService {
	return &tagService{
		container: container,
		repo:      repository.NewTagRepository(container.GetDB()),
	}
}

func (ts *tagService) CreateTag(tagDto dto.CreateOrUpdateTag) (entity.Tag, error) {
	tag := entity.Tag{
		Name:  tagDto.Name,
		Color: SetDefaultColorIfEmpty(tagDto.Color),
	}

	return ts.repo.CreateTag(tag)
}

func (ts *tagService) UpdateTag(id uint, tagDto dto.CreateOrUpdateTag) (entity.Tag, error) {
	now := time.Now()
	tag := entity.Tag{
		ID:        id,
		UpdatedAt: &now,
		Name:      tagDto.Name,
		Color:     SetDefaultColorIfEmpty(tagDto.Color),
	}

	return ts.repo.UpdateTag(id, tag)
}

func (ts *tagService) PatchTag(id uint, patchTag dto.PatchTag) (entity.Tag, error) {
	tag := entity.Tag{}

	if patchTag.Color != nil {
		tag.Color = SetDefaultColorIfEmpty(*patchTag.Color)
	}

	if patchTag.Name != nil {
		tag.Name = *patchTag.Name
	}

	return ts.repo.PatchTag(id, tag)
}

func (ts *tagService) DeleteTag(id uint) (entity.Tag, error) {
	return ts.repo.DeleteTag(id)
}

func (ts *tagService) GetTags(q dto.TagListQuery) ([]entity.Tag, error) {
	return ts.repo.GetTags(q)
}

func (ts *tagService) GetTagByID(id uint) (entity.Tag, error) {
	return ts.repo.GetTagByID(id)
}

func SetDefaultColorIfEmpty(color string) string {
	if color == "" {
		return DEFAULT_COLOR
	}

	return color
}

package service

import (
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/internal/repository"
)

type TagService interface {
	CreateTag(createTag dto.TagRequest) (entity.Tag, error)
	GetTags() ([]entity.Tag, error)
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

func (ts *tagService) CreateTag(tagDto dto.TagRequest) (entity.Tag, error) {
	return ts.repo.CreateTag(tagDto)
}

func (ts *tagService) GetTags() ([]entity.Tag, error) {
	return ts.repo.GetTags()
}

func (ts *tagService) GetTagByID(id uint) (entity.Tag, error) {
	return ts.repo.GetTagByID(id)
}

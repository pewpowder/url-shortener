package service

import (
	"time"

	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/internal/repository"
	"gorm.io/gorm"
)

func takePtr[T any](value T) *T {
	return &value
}

var links = []entity.Link{
	{
		Model:        gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		ShortCode:    "aHR0cHM6Ly9nb29nbGUuY29t",
		OriginalURL:  "https://google.com",
		ExpiresAt:    takePtr(time.Now().Add(time.Hour * 10)),
		IsActive:     true,
		IsPrivate:    false,
		PasswordHash: nil,
		UserID:       0,
		MaxClicks:    100,
		Tags:         nil,
	},
	{
		Model:        gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		ShortCode:    "aHR0cHM6Ly9nb29nbGUuY29t",
		OriginalURL:  "https://google.com",
		ExpiresAt:    takePtr(time.Now().Add(time.Hour * 5)),
		IsActive:     false,
		IsPrivate:    false,
		PasswordHash: nil,
		UserID:       0,
		MaxClicks:    0,
		Tags:         nil,
	},
	{
		Model:        gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		ShortCode:    "aHR0cHM6Ly9nb29nbGUuY29t",
		OriginalURL:  "https://google.com",
		ExpiresAt:    takePtr(time.Now().Add(time.Hour * 20)),
		IsActive:     false,
		IsPrivate:    true,
		PasswordHash: takePtr("dmVyeV9zdHJvbmdfcGFzc3dvcmRfaGFzaA=="),
		UserID:       0,
		MaxClicks:    123,
		Tags:         nil,
	},
}

type LinkService interface {
	GetLinks(query *dto.LinkListQuery) ([]entity.Link, error)
	// CreateLink(data *dto.CreateLink) entity.Link
}

type linkService struct {
	container container.Container
	repo      repository.LinkRepository
}

func NewLinkService(container container.Container) LinkService {
	return &linkService{
		container: container,
		repo:      repository.NewLinkRepository(container.GetDB()),
	}
}

func (ls *linkService) GetLinks(query *dto.LinkListQuery) ([]entity.Link, error) {
	return ls.repo.GetLinks("", "")
}

// func (ls *linkService) CreateLink(data *dto.CreateLink) entity.Link {
// }

// TODO: Imeplement method for password determination with rules (like if IsPrivate is false I should reset password to nil etc)

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/internal/repository"
	se "github.com/pewpowder/url-shortener/pkg/errors"
	"github.com/pewpowder/url-shortener/pkg/utils"
)

const MAX_ATTEMPTS = 5 // attempts to generate the hash

// var links = []entity.Link{
// 	{
// 		Model:        gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
// 		ShortCode:    "aHR0cHM6Ly9nb29nbGUuY29t",
// 		OriginalURL:  "https://google.com",
// 		ExpiresAt:    takePtr(time.Now().Add(time.Hour * 10)),
// 		IsActive:     true,
// 		IsPrivate:    false,
// 		PasswordHash: nil,
// 		UserID:       0,
// 		MaxClicks:    100,
// 		Tags:         nil,
// 	},
// 	{
// 		Model:        gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
// 		ShortCode:    "aHR0cHM6Ly9nb29nbGUuY29t",
// 		OriginalURL:  "https://google.com",
// 		ExpiresAt:    takePtr(time.Now().Add(time.Hour * 5)),
// 		IsActive:     false,
// 		IsPrivate:    false,
// 		PasswordHash: nil,
// 		UserID:       0,
// 		MaxClicks:    0,
// 		Tags:         nil,
// 	},
// 	{
// 		Model:        gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
// 		ShortCode:    "aHR0cHM6Ly9nb29nbGUuY29t",
// 		OriginalURL:  "https://google.com",
// 		ExpiresAt:    takePtr(time.Now().Add(time.Hour * 20)),
// 		IsActive:     false,
// 		IsPrivate:    true,
// 		PasswordHash: takePtr("dmVyeV9zdHJvbmdfcGFzc3dvcmRfaGFzaA=="),
// 		UserID:       0,
// 		MaxClicks:    123,
// 		Tags:         nil,
// 	},
// }

type LinkService interface {
	GetLinks(ctx context.Context, query dto.LinkListQuery) ([]entity.Link, error)
	GetLinkById(ctx context.Context, id uint) (entity.Link, error)
	GetLinkByCode(ctx context.Context, code string) (entity.Link, error)
	CreateLink(ctx context.Context, data dto.CreateLink) (entity.Link, error)
	UpdateLink(ctx context.Context, data dto.CreateLink) (entity.Link, error)
	PatchLink(ctx context.Context, data dto.PatchLink) (entity.Link, error)
	DeleteLink(ctx context.Context, id uint) error
}

type linkService struct {
	container  container.Container
	repo       repository.LinkRepository
	tagService TagService
}

func NewLinkService(container container.Container, tagService TagService) LinkService {
	return &linkService{
		container:  container,
		repo:       repository.NewLinkRepository(container.GetDB()),
		tagService: tagService,
	}
}

func (ls *linkService) GetLinks(ctx context.Context, query dto.LinkListQuery) ([]entity.Link, error) {
	return ls.repo.GetLinks(query)
}

func (ls *linkService) GetLinkById(ctx context.Context, id uint) (entity.Link, error) {
	return ls.repo.GetLinkById(id)
}

func (ls *linkService) GetLinkByCode(ctx context.Context, code string) (entity.Link, error) {
	return ls.repo.GetLinkByShortCode(code)
}

func (ls *linkService) CreateLink(ctx context.Context, createLink dto.CreateLink) (entity.Link, error) {
	shortCode, err := ls.generateUniqueCode(MAX_ATTEMPTS)
	if err != nil {
		return entity.Link{}, fmt.Errorf("failed to generate short code for %s: %w", createLink.URL, err)
	}

	isActive := utils.DerefBool(createLink.IsActive, true)
	isPrivate := utils.DerefBool(createLink.IsPrivate, false)

	expiresAt, err := utils.ParseOptionalTime(createLink.ExpiresAt)
	if err != nil {
		return entity.Link{}, se.NewServiceError("invalid expires_at date", se.ErrCodeBadRequest, se.GetCodeTextByCode(se.ErrCodeBadRequest), err)
	}

	// Calculate hash only if link private and password provided
	var passHash *string
	if isPrivate && createLink.Password != "" {
		h, err := utils.CreateHash(createLink.Password)
		if err != nil {
			return entity.Link{}, se.NewServiceError(
				"failed to generate password hash",
				se.ErrCodeInternal,
				se.GetCodeTextByCode(se.ErrCodeInternal),
				err,
			)
		}

		passHash = &h
	}

	linkTags := make([]entity.Tag, 0, len(createLink.Tags))
	if len(createLink.Tags) > 0 {
		q := dto.TagListQuery{
			Pagination: dto.Pagination{
				Page: 1,
				Size: len(createLink.Tags),
			},
			Sort:   dto.Sort{},
			Filter: &dto.TagFilter{Names: createLink.Tags},
		}

		tags, err := ls.tagService.GetTags(q)
		if err != nil {
			return entity.Link{}, se.NewServiceError(
				"failed to get tags",
				se.ErrCodeInternal,
				se.GetCodeTextByCode(se.ErrCodeInternal),
				err,
			)
		}

		if (len(createLink.Tags) - len(tags)) != 0 {
			tagsNames := make([]string, 0, len(tags))
			for _, tag := range tags {
				tagsNames = append(tagsNames, tag.Name)
			}

			nonExistedTags := utils.Difference(createLink.Tags, tagsNames)
			return entity.Link{}, se.NewServiceError(
				fmt.Sprintf("Not found following tags: %v", nonExistedTags),
				se.ErrCodeBadRequest,
				se.GetCodeTextByCode(se.ErrCodeBadRequest),
				errors.New("not found tags"),
			)
		}

		linkTags = tags
	}

	// TODO: FINISH CREATE LINK METHOD (GENERATE PASSWORD PROPERLY)

	link := entity.Link{
		ShortCode:    shortCode,
		OriginalURL:  createLink.URL,
		IsPrivate:    isPrivate,
		IsActive:     isActive,
		PasswordHash: passHash,
		MaxClicks:    createLink.MaxClicks,
		Tags:         linkTags,
		ExpiresAt:    expiresAt,
	}

	return ls.repo.CreateLink(link)
}

func (ls *linkService) UpdateLink(ctx context.Context, data dto.CreateLink) (entity.Link, error) {
	return entity.Link{}, se.NewServiceError("not implemented", se.ErrCodeInternal, se.GetCodeTextByCode(se.ErrCodeInternal), errors.New("update not implemented"))
}

func (ls *linkService) PatchLink(ctx context.Context, data dto.PatchLink) (entity.Link, error) {
	return entity.Link{}, se.NewServiceError("not implemented", se.ErrCodeInternal, se.GetCodeTextByCode(se.ErrCodeInternal), errors.New("update not implemented"))
}

func (ls *linkService) DeleteLink(ctx context.Context, id uint) error {
	return ls.repo.DeleteLink(uint(id))
}

func (ls *linkService) generateUniqueCode(maxAttempts int) (string, error) {
	for range maxAttempts {
		code := utils.GenerateNanoid()
		_, err := ls.repo.GetLinkByShortCode(code)

		if err != nil {
			var serviceErr *se.ServiceError
			if errors.As(err, &serviceErr) && serviceErr.Code == se.ErrCodeNotFound {
				return code, nil
			}

			return "", err
		}
	}

	return "", se.NewServiceError(
		fmt.Sprintf("failed to generate unique short code after %d attempts", maxAttempts),
		se.ErrCodeInternal,
		se.GetCodeTextByCode(se.ErrCodeInternal),
		errors.New("unique generator error"),
	)
}

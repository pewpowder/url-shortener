package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/internal/repository"
	se "github.com/pewpowder/url-shortener/pkg/errors"
)

// func takePtr[T any](value T) *T {
// 	return &value
// }

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
	CreateLink(ctx context.Context, data dto.CreateLink) (entity.Link, error)
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

func (ls *linkService) GetLinks(ctx context.Context, query dto.LinkListQuery) ([]entity.Link, error) {
	// if userIDStr, err := utils.GetUserID(ctx); err == nil {
	// 	userID, _ := strconv.ParseUint(userIDStr, 10, 64)
	// 	query.UserID = userID.(uint64)
	// }

	return ls.repo.GetLinks(query)
}

func (ls *linkService) CreateLink(ctx context.Context, data dto.CreateLink) (entity.Link, error) {
	// userID, _ := utils.GetUserID(ctx)
	// shortCode, err := ls.GenerateUniqueShortCode(data.URL, 8, 10)
	// if err != nil {
	// 	return entity.Link{}, err
	// }

	return entity.Link{}, nil
}

func (ls *linkService) GenerateUniqueShortCode(url string, length int, maxAttempts int) (string, error) {
	// code, err := GenerateShortCode(url, length, nil)

	// if err != nil {
	// 	return "", se.NewServiceError("failed to generate short code", se.ErrTypeInternal, err)
	// }

	// if link, err := ls.repo.GetLinkByShortCode(code); err == nil {
	// 	return link.ShortCode, nil
	// }

	var salt []byte
	for attempt := range maxAttempts {
		if attempt > 0 {
			salt = make([]byte, 4)
			if _, err := rand.Read(salt); err != nil {
				return "", err
			}
		}

		candidateCode, err := GenerateShortCode(url, length, salt)
		if err != nil {
			return "", err
		}

		// Check for collision
		link, err := ls.repo.GetLinkByShortCode(candidateCode)
		switch {
		case err == nil && link.OriginalURL == url:
			return candidateCode, nil // Код уже существует для этого URL
		case err == nil:
			continue // Collision with different URL
		case repository.IsGormNotFoundError(err):
			// Code available, save to DB
			// if err := s.db.SaveURL(code, url); err == nil {
			// 	return code, nil
			// } else if s.db.IsDuplicateError(err) {
			// 	continue // Collision from concurrent insert
			// } else {
			// 	return "", err
			// }
		default:
			return "", err
		}
	}

	return "", &se.ServiceError{
		Message: fmt.Sprintf("failed to generate unique code after %d attempts", maxAttempts),
		Type:    se.ErrTypeInternal,
		Err:     errors.New("failed to generate unique code"),
	}
}

// Deterministic short code generation
func GenerateShortCode(url string, length int, salt []byte) (string, error) {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if length < 1 || length > 11 {
		return "", &se.ServiceError{}
	}

	// Compute SHA-256 hash of URL + salt
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hasher.Write(salt)
	hash := hasher.Sum(nil)

	// Use first 8 bytes to create a uint64
	num := binary.BigEndian.Uint64(hash[:8])

	// For lengths 1-10, reduce number space to base^length
	if length <= 10 {
		basePow := uint64(1)
		for range length {
			basePow *= 62
		}
		num %= basePow
	}

	// Convert to base62
	code := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		remainder := num % 62
		num /= 62
		code[i] = charset[remainder]
	}

	return string(code), nil
}

// TODO: how to correctly handle first attempt for password? Should I generate salt for first attempt? Or just use empty salt?
func GeneratePasswordHash(password string) string {
	return ""
}

// Get password hash by curtain rules
// 1. IsPrivate is true -> password is not nil
// 2. IsPrivate is false -> reset password to nil (if password is not empty)
// 3. IsPrivate is empty -> don't change password and return it
func GetPasswordHash(password string, isPrivate bool) string {
	return ""
}

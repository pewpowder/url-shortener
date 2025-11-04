package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jaevor/go-nanoid"
	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/pkg/logger"
)

func DerefBool(b *bool, defaultVal bool) bool {
	if b == nil {
		return defaultVal
	}
	return *b
}

func ParseOptionalTime(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	parsedTime, err := time.Parse(config.TIME_FORMAT, *value)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %w", err)
	}

	// Normalize to UTC and ensure nanosecond precision is stripped
	utcTime := parsedTime.UTC().Round(time.Second)
	return &utcTime, nil
}

func ParseTime(value *string) (time.Time, error) {
	parsedTime, err := time.Parse(config.TIME_FORMAT, *value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %w", err)
	}

	return parsedTime, nil
}

// Gin default doesn't split comma-separated strings into slices (param=str1,str2,str3... will not be split into slices)
func GinSplitString(data []string) []string {
	if len(data) == 1 {
		return strings.Split(data[0], ",")
	}

	return data
}

func CreateHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CompareHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func GenerateNanoid() string {
	nanoidGen, err := nanoid.Canonic()

	if err != nil {
		logger.Get().Fatal().Msg("Can't initiate nanoid generator")
	}

	return nanoidGen()
}

func Difference[T comparable](a, b []T) []T {
	lookup := make(map[T]struct{}, len(b))
	for _, val := range b {
		lookup[val] = struct{}{}
	}

	// Collect elements from a not in b
	var diff []T
	for _, val := range a {
		if _, found := lookup[val]; !found {
			diff = append(diff, val)
		}
	}

	return diff
}

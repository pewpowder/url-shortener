package utils

import (
	"fmt"
	"time"

	"github.com/pewpowder/url-shortener/internal/config"
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

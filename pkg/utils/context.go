package utils

import (
	"context"
	"errors"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// SetUserID добавляет USER_ID в context
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID извлекает USER_ID из context
func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", errors.New("user_id not found in context")
	}
	return userID, nil
}

// MustGetUserID извлекает USER_ID из context или паникует
func MustGetUserID(ctx context.Context) string {
	userID, err := GetUserID(ctx)
	if err != nil {
		panic("user_id is required but not found in context")
	}
	return userID
}

// HasUserID проверяет наличие USER_ID в context
func HasUserID(ctx context.Context) bool {
	_, err := GetUserID(ctx)
	return err == nil
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	se "github.com/pewpowder/url-shortener/pkg/errors"
)

type AuthService interface {
	GetUserIDByToken(ctx context.Context, token string) (string, error)
}

type authService struct {
	authServiceURL string
	httpClient     *http.Client
}

type AuthResponse struct {
	UserID string `json:"user_id"`
}

func NewAuthService(authServiceURL string) AuthService {
	return &authService{
		authServiceURL: authServiceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (as *authService) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", as.authServiceURL, nil)
	if err != nil {
		return "", se.NewServiceError("failed to create request", se.ErrCodeInternal, se.ErrInternal, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := as.httpClient.Do(req)
	if err != nil {
		return "", se.NewServiceError("failed to make request to auth service", se.ErrCodeInternal, se.ErrInternal, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", se.NewServiceError(fmt.Sprintf("auth service returned status %d", resp.StatusCode), se.ErrCodeInternal, se.ErrInternal, nil)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", se.NewServiceError("failed to decode auth response", se.ErrCodeInternal, se.ErrInternal, err)
	}

	if authResp.UserID == "" {
		return "", se.NewServiceError("invalid token or user not found", se.ErrCodeUnauthorized, se.ErrUnauthorized, nil)
	}

	return authResp.UserID, nil
}

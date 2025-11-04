package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/pkg/logger"
)

func HandleError(c *gin.Context, err error) {
	logger.Get().Err(err)

	var se *ServiceError
	if errors.As(err, &se) {
		c.JSON(ErrorCodeToHTTPStatus(se.Code), gin.H{
			"code":    ErrorCodeToHTTPStatus(se.Code),
			"error":   se.CodeText,
			"details": se.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"error":   "Internal Server Error",
		"details": "Internal Server Error",
	})
}

func InvalidQueryParams(c *gin.Context, err error) {
	HandleError(c, NewServiceError(fmt.Sprintf("invalid query params: %s", err), ErrCodeBadRequest, ErrInvalidRequest, err))
}

func InvalidRequestData(c *gin.Context, err error) {
	HandleError(c, NewServiceError(fmt.Sprintf("invalid request data: %s", err), ErrCodeBadRequest, ErrInvalidRequest, err))
}

func UnauthorizedError(c *gin.Context, err error) {
	HandleError(c, NewServiceError(fmt.Sprintf("unauthorized: %s", err), ErrCodeUnauthorized, ErrUnauthorized, err))
}

func InternalServerError(c *gin.Context, err error) {
	HandleError(c, NewServiceError(fmt.Sprintf("internal server error: %s", err), ErrCodeInternal, ErrInternal, err))
}

func ErrorCodeToHTTPStatus(errorType ErrorCode) int {
	switch errorType {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeExpired:
		return http.StatusGone
	default:
		return http.StatusInternalServerError
	}
}

func GetErrorByErrorCode(code ErrorCode) string {
	switch code {
	case ErrCodeBadRequest:
		return "Bad Request"
	case ErrCodeNotFound:
		return "Not Found"
	case ErrCodeConflict:
		return "Conflict"
	case ErrCodeUnauthorized:
		return "Unauthorized"
	case ErrCodeForbidden:
		return "Forbidden"
	case ErrCodeExpired:
		return "Gone"
	default:
		return "Internal Server Error"
	}
}

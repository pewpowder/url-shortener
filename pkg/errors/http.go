package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	var se *ServiceError
	if errors.As(err, &se) {
		c.JSON(ErrorTypeToHTTPStatus(se.Type), gin.H{"error": se.Message})
		return
	}

	// response with default error message because message in err can contain secret information
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

func ErrorTypeToHTTPStatus(errorType ErrorType) int {
	switch errorType {
	case ErrTypeValidation:
		return http.StatusBadRequest
	case ErrTypeNotFound:
		return http.StatusNotFound
	case ErrTypeConflict:
		return http.StatusConflict
	case ErrTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrTypeForbidden:
		return http.StatusForbidden
	case ErrTypeExpired:
		return http.StatusGone
	default:
		return http.StatusInternalServerError
	}
}

func InvalidQueryParams(c *gin.Context, err error) {
	HandleError(c, NewServiceError(fmt.Sprintf("invalid query params: %s", err), ErrTypeValidation, err))
}

func InvalidRequestData(c *gin.Context, err error) {
	HandleError(c, NewServiceError(fmt.Sprintf("invalid request data: %s", err), ErrTypeValidation, err))
}

func UnauthorizedError(c *gin.Context, err error) {
	HandleError(c, NewServiceError(fmt.Sprintf("unauthorized: %s", err), ErrTypeUnauthorized, err))
}

func InternalServerError(c *gin.Context, err error) {
	HandleError(c, NewServiceError(fmt.Sprintf("internal server error: %s", err), ErrTypeInternal, err))
}

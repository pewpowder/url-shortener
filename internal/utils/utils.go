package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	se "github.com/pewpowder/url-shortener/pkg/errors"
)

func HandleError(c *gin.Context, err error) {
	var se *se.ServiceError
	if errors.As(err, &se) {
		c.JSON(ErrorTypeToHTTPStatus(se.Type), gin.H{"error": se.Message})
		return
	}

	// response with default error message because message in err can contain secret information
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

func ErrorTypeToHTTPStatus(errorType se.ErrorType) int {
	switch errorType {
	case se.ErrTypeValidation:
		return http.StatusBadRequest
	case se.ErrTypeNotFound:
		return http.StatusNotFound
	case se.ErrTypeConflict:
		return http.StatusConflict
	case se.ErrTypeUnauthorized:
		return http.StatusUnauthorized
	case se.ErrTypeForbidden:
		return http.StatusForbidden
	case se.ErrTypeExpired:
		return http.StatusGone
	default:
		return http.StatusInternalServerError
	}
}

func InvalidQueryParams(c *gin.Context, err error) {
	HandleError(c, se.NewServiceError(fmt.Sprintf("invalid query params: %s", err), se.ErrTypeValidation, err))
}

func InvalidRequestData(c *gin.Context, err error) {
	HandleError(c, se.NewServiceError(fmt.Sprintf("invalid request data: %s", err), se.ErrTypeValidation, err))
}

func InternalServerError(c *gin.Context, err error) {
	HandleError(c, se.NewServiceError(fmt.Sprintf("internal server error: %s", err), se.ErrTypeInternal, err))
}

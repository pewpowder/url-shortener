package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestError struct {
	Status  int    `json:"status_code"`
	Message string `json:"message"`
}

func (rq *RequestError) Error() string {
	return rq.Message
}

func NewRequestError(status int, message string) error {
	return &RequestError{
		Status:  status,
		Message: message,
	}
}

func InvalidQueryParams(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid query params: %s", err)})
}

func InvalidBody(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid body: %s", err)})
}

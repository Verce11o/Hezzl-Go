package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

func MapHTTPError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		return http.StatusBadRequest, "bad request"
	}
	return http.StatusInternalServerError, "server error"

}

func WithHTTPError(c *gin.Context, err error) {
	status, message := MapHTTPError(err)
	c.JSON(status, gin.H{
		"message": message,
	})
	return
}

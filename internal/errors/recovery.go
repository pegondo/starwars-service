package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware is a middleware that reads the gin context errors and
// transforms them into a JSON response. If the error is a ResponseError, the
// response will have its content; if not, it will be an INTERNAL_SERVER_ERROR.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if err := c.Errors.Last(); err != nil {
			statusCode := http.StatusInternalServerError
			jsonErr := InternalServerError()
			var respErr *ResponseError
			if errors.As(err, &respErr) {
				statusCode = c.Writer.Status()
				jsonErr = respErr
			}
			c.JSON(statusCode, jsonErr)
		}
	}
}

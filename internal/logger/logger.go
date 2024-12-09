package logger

import (
	"errors"
	"starwars/service/internal/request"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// loggerKey is the key to find the logger in the request context.
var loggerKey = "logger"

// Middleware is a middleware that attaches a logger to the request
// context.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := request.RequestId(c)
		logger := log.With().Str("reqId", reqId).Logger()
		c.Set(loggerKey, logger)

		c.Next()
	}
}

// Logger returns the logger attached to the given request context.
func Logger(c *gin.Context) (l zerolog.Logger, err error) {
	logger, exists := c.Get(loggerKey)
	if !exists {
		return l, errors.New("error while retrieving the logger :: no logger in the request")
	}

	l, ok := logger.(zerolog.Logger)
	if !ok {
		return l, errors.New("error while retrieving the logger :: invalid logger instance type")
	}

	return l, nil
}

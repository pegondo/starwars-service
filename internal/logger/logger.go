package logger

import (
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

// Logger returns the logger attached to the given request context. If there is
// no log in the request context, Logger logs an error and returns the default
// one.
func Logger(c *gin.Context) (l zerolog.Logger) {
	logger, exists := c.Get(loggerKey)
	if !exists {
		log.Error().Msg("no logger in the request")
		return l
	}

	l, ok := logger.(zerolog.Logger)
	if !ok {
		log.Error().Msg("invalid logger instance type")
	}
	return l
}

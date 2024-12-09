package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// requestIdKey is the key to find the request id in the request context.
const requestIdKey = "request-id"

// RequestIdMiddleware is a middleware that attaches a request id to the
// requests.
func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := uuid.NewString()
		c.Set(requestIdKey, reqId)

		c.Next()
	}
}

// RequestId returns the request id linked to the given request context. If
// the context has no request id, RequestId returns "".
func RequestId(c *gin.Context) string {
	reqId, exists := c.Get(requestIdKey)
	if !exists {
		return ""
	}
	reqIdStr, ok := reqId.(string)
	if !ok {
		return ""
	}
	return reqIdStr
}

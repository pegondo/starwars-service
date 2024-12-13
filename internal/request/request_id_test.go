package request_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pegondo/starwars/service/internal/request"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func buildReqIdRouter(handler gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Use(request.RequestIdMiddleware())
	r.GET("/", handler)
	return r
}

func TestRequestId(t *testing.T) {
	var reqId string
	handler := func(c *gin.Context) {
		reqId = request.RequestId(c)
	}
	r := buildReqIdRouter(handler)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.NotEmpty(t, reqId)
}

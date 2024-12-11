package logger_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"starwars/service/internal/logger"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func buildRouter(handler gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Use(logger.Middleware())
	r.GET("/", handler)
	return r
}

func TestLogger(t *testing.T) {
	lw := httptest.NewRecorder()
	log.Logger = zerolog.New(lw)

	handler := func(c *gin.Context) {
		l := logger.Logger(c)
		l.Info().Msg("<msg>")
		c.Status(http.StatusNoContent)
	}
	r := buildRouter(handler)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the logger printed the request id.
	var lwBodyStruct struct {
		ReqId *string `json:"reqId"`
	}
	json.Unmarshal(lw.Body.Bytes(), &lwBodyStruct)
	require.NotNil(t, lwBodyStruct.ReqId)
}

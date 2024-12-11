package errors_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"starwars/service/internal/errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func buildRouter(handler gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Use(errors.RecoveryMiddleware())
	r.GET("/", handler)
	return r
}

func TestRecoveryMiddleware_NoError(t *testing.T) {
	noErrHandler := func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	}
	r := buildRouter(noErrHandler)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestRecoveryMiddleware_ResponseError(t *testing.T) {
	responseErr := &errors.ResponseError{
		ErrorCode:    "<error-code>",
		ErrorMessage: "<error-message>",
	}
	responseErrHandler := func(c *gin.Context) {
		c.AbortWithError(http.StatusBadRequest, responseErr)
	}
	r := buildRouter(responseErrHandler)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	responseErrJson, err := json.Marshal(responseErr)
	require.NoError(t, err)
	responseErrStr := string(responseErrJson)
	require.Equal(t, responseErrStr, w.Body.String())
}

func TestRecoveryMiddleware_NotResponseError(t *testing.T) {
	errHandler := func(c *gin.Context) {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("<error>"))
	}
	r := buildRouter(errHandler)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	responseErrJson, err := json.Marshal(errors.InternalServerError())
	require.NoError(t, err)
	responseErrStr := string(responseErrJson)
	require.Equal(t, responseErrStr, w.Body.String())
}

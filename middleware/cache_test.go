package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFrontendCachePolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		requestPath string
		want        string
	}{
		{requestPath: "/", want: "no-store, no-cache, must-revalidate"},
		{requestPath: "/console", want: "no-store, no-cache, must-revalidate"},
		{requestPath: "/dashboard/overview", want: "no-store, no-cache, must-revalidate"},
		{requestPath: "/index.html", want: "no-store, no-cache, must-revalidate"},
		{requestPath: "/static/js/index.abc123.js", want: "public, max-age=31536000, immutable"},
		{requestPath: "/assets/index-abc123.js", want: "public, max-age=31536000, immutable"},
		{requestPath: "/logo.png", want: "public, max-age=604800"},
	}

	for _, test := range tests {
		t.Run(test.requestPath, func(t *testing.T) {
			engine := gin.New()
			engine.Use(Cache())
			engine.GET(test.requestPath, func(c *gin.Context) { c.Status(http.StatusOK) })
			request := httptest.NewRequest(http.MethodGet, test.requestPath, nil)
			response := httptest.NewRecorder()

			engine.ServeHTTP(response, request)

			assert.Equal(t, test.want, response.Header().Get("Cache-Control"))
			assert.Equal(t, frontendCacheVersion, response.Header().Get("Cache-Version"))
		})
	}
}

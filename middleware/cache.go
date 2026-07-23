package middleware

import (
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

const frontendCacheVersion = "20260722-spa-cache-recovery-v1"

func Cache() func(c *gin.Context) {
	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		extension := strings.ToLower(path.Ext(requestPath))
		switch {
		case requestPath == "/", extension == "", extension == ".html":
			// SPA documents contain build-specific chunk names. They must never be
			// reused after a deployment, including deep links such as /console.
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			c.Header("CDN-Cache-Control", "no-store")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		case strings.HasPrefix(requestPath, "/static/"), strings.HasPrefix(requestPath, "/assets/"):
			// Both Rsbuild and Vite emit content-hashed assets.
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		default:
			c.Header("Cache-Control", "public, max-age=604800")
		}
		c.Header("Cache-Version", frontendCacheVersion)
		c.Next()
	}
}

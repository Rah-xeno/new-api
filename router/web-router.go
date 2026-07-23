package router

import (
	"embed"
	"net/http"
	"path"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const frontendAssetReloadScript = `(function () {
  var key = 'newapi_asset_reload_at';
  var now = Date.now();
  var shouldReload = true;
  try {
    var last = Number(sessionStorage.getItem(key) || 0);
    shouldReload = !last || now - last > 30000;
    if (shouldReload) sessionStorage.setItem(key, String(now));
  } catch (_) {}
  if (shouldReload) {
    var url = new URL(window.location.href);
    url.searchParams.set('__newapi_asset_reload', String(now));
    window.location.replace(url.toString());
  }
})();
`

func isMissingFrontendJavaScript(requestPath string) bool {
	if !strings.HasPrefix(requestPath, "/static/js/") && !strings.HasPrefix(requestPath, "/assets/") {
		return false
	}
	extension := strings.ToLower(path.Ext(requestPath))
	return extension == ".js" || extension == ".mjs"
}

// ThemeAssets holds the embedded frontend assets.
type ThemeAssets struct {
	BuildFS   embed.FS
	IndexPage []byte
}

func SetWebRouter(router *gin.Engine, assets ThemeAssets) {
	webFS := common.EmbedFolder(assets.BuildFS, "web/default/dist")

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())
	router.Use(middleware.Cache())
	router.Use(static.Serve("/", webFS))
	router.NoRoute(func(c *gin.Context) {
		c.Set(middleware.RouteTagKey, "web")
		requestPath := c.Request.URL.Path
		if isMissingFrontendJavaScript(requestPath) {
			// A cached document can request chunks removed by a newer deployment.
			// This script works as both a classic script (Rsbuild) and an ES module
			// (Vite), and forces one cache-busted document reload.
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			c.Header("CDN-Cache-Control", "no-store")
			c.Data(http.StatusOK, "application/javascript; charset=utf-8", []byte(frontendAssetReloadScript))
			return
		}
		if strings.HasPrefix(requestPath, "/v1") ||
			strings.HasPrefix(requestPath, "/api") ||
			strings.HasPrefix(requestPath, "/assets") ||
			strings.HasPrefix(requestPath, "/static") {
			c.Header("Cache-Control", "no-store")
			controller.RelayNotFound(c)
			return
		}
		if path.Ext(requestPath) != "" {
			c.Header("Cache-Control", "no-store")
			controller.RelayNotFound(c)
			return
		}
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("CDN-Cache-Control", "no-store")
		c.Data(http.StatusOK, "text/html; charset=utf-8", assets.IndexPage)
	})
}

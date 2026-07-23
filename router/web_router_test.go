package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingFrontendJavaScriptDetection(t *testing.T) {
	assert.True(t, isMissingFrontendJavaScript("/static/js/index.old.js"))
	assert.True(t, isMissingFrontendJavaScript("/assets/dashboard-old.MJS"))
	assert.False(t, isMissingFrontendJavaScript("/static/css/index.old.css"))
	assert.False(t, isMissingFrontendJavaScript("/api/static/js/index.js"))
	assert.Contains(t, frontendAssetReloadScript, "window.location.replace")
}

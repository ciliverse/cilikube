package initialization

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterDesktopStatic serves the React SPA from webRoot when present.
// API routes must be registered first; this attaches NoRoute fallback last.
func RegisterDesktopStatic(router *gin.Engine, webRoot string) {
	webRoot = strings.TrimSpace(webRoot)
	if webRoot == "" {
		return
	}
	index := filepath.Join(webRoot, "index.html")
	if st, err := os.Stat(index); err != nil || st.IsDir() {
		return
	}

	fileServer := http.FileServer(http.Dir(webRoot))
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") ||
			path == "/health" || path == "/ready" || path == "/live" ||
			path == "/metrics" || path == "/version" ||
			strings.HasPrefix(path, "/uploads/") {
			c.Status(http.StatusNotFound)
			return
		}

		// Prefer real files (assets/*); otherwise SPA index.html
		clean := path
		if clean == "/" || clean == "" {
			c.File(index)
			return
		}
		full := filepath.Join(webRoot, filepath.Clean("/"+clean))
		if !strings.HasPrefix(full, filepath.Clean(webRoot)) {
			c.Status(http.StatusNotFound)
			return
		}
		if st, err := os.Stat(full); err == nil && !st.IsDir() {
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}
		c.File(index)
	})
}

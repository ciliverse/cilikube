package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

type PluginHandler struct {
	root string
}

func NewPluginHandler(root string) *PluginHandler {
	if root == "" {
		root = "plugins"
	}
	return &PluginHandler{root: root}
}

func (h *PluginHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": service.LoadPlugins(h.root)})
}

func (h *PluginHandler) Asset(c *gin.Context) {
	id := c.Param("id")
	name := strings.TrimPrefix(c.Param("filepath"), "/")
	if name == "" {
		name = "index.html"
	}
	if strings.Contains(id, "..") || strings.Contains(name, "..") {
		c.Status(http.StatusBadRequest)
		return
	}
	path := filepath.Join(h.root, id, filepath.Clean(name))
	if _, err := os.Stat(path); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.File(path)
}

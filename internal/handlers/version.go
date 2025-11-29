package handlers

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

// These variables will be set at compile time via ldflags
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// VersionInfo version information structure
type VersionInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
	Arch      string `json:"arch"`
}

// GetVersion retrieves version information
func GetVersion(c *gin.Context) {
	versionInfo := VersionInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    versionInfo,
		"message": "success",
	})
}

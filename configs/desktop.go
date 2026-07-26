package configs

import (
	"net"
	"os"
	"strings"
)

// IsDesktop reports whether the process is running as a desktop sidecar.
func IsDesktop() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("CILIKUBE_DESKTOP")))
	return v == "1" || v == "true" || v == "yes"
}

// ApplyDesktopOverrides applies env-driven desktop defaults after config load.
// Safe to call when not in desktop mode (only CILIKUBE_ADDR / CILIKUBE_WEB_ROOT apply).
func ApplyDesktopOverrides(cfg *Config) {
	if cfg == nil {
		return
	}

	if addr := strings.TrimSpace(os.Getenv("CILIKUBE_ADDR")); addr != "" {
		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			if host != "" {
				cfg.Server.Host = host
			}
			if port != "" {
				cfg.Server.Port = port
			}
		}
	}

	if web := strings.TrimSpace(os.Getenv("CILIKUBE_WEB_ROOT")); web != "" {
		cfg.Server.WebRoot = web
	}

	if !IsDesktop() {
		return
	}

	// Desktop: loopback only unless ADDR already set a host.
	if strings.TrimSpace(cfg.Server.Host) == "" {
		cfg.Server.Host = "127.0.0.1"
	}
	if strings.TrimSpace(cfg.Server.Port) == "" {
		cfg.Server.Port = "17880"
	}
	// Prefer release logging for shipped builds.
	if cfg.Server.Mode == "" || cfg.Server.Mode == "debug" {
		if os.Getenv("CILIKUBE_SERVER_MODE") == "" {
			cfg.Server.Mode = "release"
		}
	}
	// Do not use public OAuth in the local desktop app by default.
	cfg.OAuth.GitHub.Enabled = false
}

package configs

import (
	"testing"
)

func TestApplyDesktopOverridesAddrAndWebRoot(t *testing.T) {
	t.Setenv("CILIKUBE_DESKTOP", "1")
	t.Setenv("CILIKUBE_ADDR", "127.0.0.1:19090")
	t.Setenv("CILIKUBE_WEB_ROOT", "/tmp/cilikube-web")

	cfg := &Config{}
	cfg.Server.Port = "8080"
	cfg.OAuth.GitHub.Enabled = true

	ApplyDesktopOverrides(cfg)

	if cfg.Server.Host != "127.0.0.1" {
		t.Fatalf("host=%q", cfg.Server.Host)
	}
	if cfg.Server.Port != "19090" {
		t.Fatalf("port=%q", cfg.Server.Port)
	}
	if cfg.Server.WebRoot != "/tmp/cilikube-web" {
		t.Fatalf("webRoot=%q", cfg.Server.WebRoot)
	}
	if cfg.OAuth.GitHub.Enabled {
		t.Fatal("expected OAuth disabled in desktop mode")
	}
}

func TestApplyDesktopOverridesNonDesktopKeepsOAuth(t *testing.T) {
	t.Setenv("CILIKUBE_DESKTOP", "")
	t.Setenv("CILIKUBE_ADDR", "127.0.0.1:19191")

	cfg := &Config{}
	cfg.OAuth.GitHub.Enabled = true
	ApplyDesktopOverrides(cfg)

	if cfg.Server.Port != "19191" {
		t.Fatalf("port=%q", cfg.Server.Port)
	}
	if !cfg.OAuth.GitHub.Enabled {
		t.Fatal("non-desktop should not force OAuth off")
	}
}

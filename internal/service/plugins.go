package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// PluginManifest is a folder under plugins/<id>/plugin.json.
type PluginManifest struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	TitleZh   string   `json:"titleZh"`
	Entry     string   `json:"entry"`
	Resources []string `json:"resources"`
	// Enabled is omitted or true by default so existing manifests stay visible.
	Enabled *bool `json:"enabled,omitempty"`
}

// LoadPlugins reads every plugins/*/plugin.json. Missing directory is empty, not an error.
func LoadPlugins(root string) []PluginManifest {
	if root == "" {
		root = "plugins"
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	var out []PluginManifest
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(root, entry.Name(), "plugin.json"))
		if err != nil {
			continue
		}
		var manifest PluginManifest
		if err := json.Unmarshal(raw, &manifest); err != nil {
			continue
		}
		if manifest.ID == "" {
			manifest.ID = entry.Name()
		}
		if strings.Contains(manifest.ID, "..") || strings.Contains(manifest.ID, "/") {
			continue
		}
		if manifest.Entry == "" {
			manifest.Entry = "index.html"
		}
		if manifest.Name == "" {
			manifest.Name = manifest.ID
		}
		if manifest.Enabled != nil && !*manifest.Enabled {
			continue
		}
		out = append(out, manifest)
	}
	return out
}

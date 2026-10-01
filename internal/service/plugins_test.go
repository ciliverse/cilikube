package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPluginsSkipsDisabled(t *testing.T) {
	dir := t.TempDir()
	write := func(id, body string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(dir, id), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, id, "plugin.json"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("on", `{"id":"on","name":"On"}`)
	write("off", `{"id":"off","name":"Off","enabled":false}`)
	got := LoadPlugins(dir)
	if len(got) != 1 || got[0].ID != "on" {
		t.Fatalf("plugins %+v", got)
	}
}

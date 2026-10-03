package handlers

import (
	"strings"
	"testing"
)

func TestRewriteRootPaths(t *testing.T) {
	in := `<a href="/app">x</a><script src="//cdn.example/a.js"></script><img src='/logo.png'>`
	got := rewriteRootPaths(in, "/proxy")
	if !strings.Contains(got, `href="/proxy/app"`) {
		t.Fatalf("href: %s", got)
	}
	if !strings.Contains(got, `src="//cdn.example/a.js"`) {
		t.Fatalf("protocol-relative changed: %s", got)
	}
	if !strings.Contains(got, `src='/proxy/logo.png'`) {
		t.Fatalf("src: %s", got)
	}
	css := rewriteRootPaths(`url(/bg.png)`, "/proxy")
	if css != `url(/proxy/bg.png)` {
		t.Fatalf("css %s", css)
	}
	html := injectBase("<head><title>t</title>", "/proxy")
	if !strings.Contains(html, `<head><base href="/proxy/">`) {
		t.Fatalf("base %s", html)
	}
}

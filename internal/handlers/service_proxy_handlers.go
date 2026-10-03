package handlers

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"
)

var (
	rootAttr = regexp.MustCompile(`(?i)(\s(?:href|src|action|poster)\s*=\s*)(["'])(/[^"']*)`)
	rootCSS  = regexp.MustCompile(`(?i)url\(\s*(["']?)(/[^)"']*)`)
)

// rewriteRootPaths points site-root URLs at the service proxy prefix.
// A leading "//" stays a protocol-relative URL.
func rewriteRootPaths(body, prefix string) string {
	if prefix == "" {
		return body
	}
	body = rootAttr.ReplaceAllStringFunc(body, func(match string) string {
		parts := rootAttr.FindStringSubmatch(match)
		if len(parts) != 4 || strings.HasPrefix(parts[3], "//") {
			return match
		}
		return parts[1] + parts[2] + prefix + parts[3]
	})
	return rootCSS.ReplaceAllStringFunc(body, func(match string) string {
		parts := rootCSS.FindStringSubmatch(match)
		if len(parts) != 3 || strings.HasPrefix(parts[2], "//") {
			return match
		}
		return "url(" + parts[1] + prefix + parts[2]
	})
}

func injectBase(html, prefix string) string {
	tag := `<base href="` + prefix + `/">`
	lower := strings.ToLower(html)
	if i := strings.Index(lower, "<head>"); i >= 0 {
		return html[:i+6] + tag + html[i+6:]
	}
	return tag + html
}

// ServiceProxyHandler forwards browser traffic to a Service through the API server proxy.
type ServiceProxyHandler struct {
	clusterManager *k8s.ClusterManager
}

func NewServiceProxyHandler(cm *k8s.ClusterManager) *ServiceProxyHandler {
	return &ServiceProxyHandler{clusterManager: cm}
}

func (h *ServiceProxyHandler) Proxy(c *gin.Context) {
	if c.Query("clusterId") == "" {
		if id, err := c.Cookie("cilikube_proxy_cluster"); err == nil && id != "" {
			q := c.Request.URL.Query()
			q.Set("clusterId", id)
			c.Request.URL.RawQuery = q.Encode()
		}
	}
	client, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		return
	}
	namespace := c.Param("namespace")
	name := c.Param("name")
	port := c.Query("port")
	if port == "" {
		if p, err := c.Cookie("cilikube_proxy_port"); err == nil && p != "" {
			port = p
		}
	}
	if port == "" {
		port = "80"
	}
	if _, err := strconv.Atoi(port); err != nil {
		c.String(http.StatusBadRequest, "port must be numeric")
		return
	}
	if k8s.IsShowcaseConfig(client.Config) {
		prefix := "/api/v1/namespaces/" + url.PathEscape(namespace) + "/services/" + url.PathEscape(name) + "/proxy"
		page := showcaseServiceHTML(namespace, name, port)
		page = injectBase(rewriteRootPaths(page, prefix), prefix)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
		return
	}
	restPath := strings.TrimPrefix(c.Param("filepath"), "/")
	if strings.Contains(restPath, "..") {
		c.String(http.StatusBadRequest, "invalid path")
		return
	}

	forwardQuery := c.Request.URL.Query()
	forwardQuery.Del("token")
	forwardQuery.Del("access_token")
	forwardQuery.Del("clusterId")
	forwardQuery.Del("port")
	target := strings.TrimRight(client.Config.Host, "/") +
		"/api/v1/namespaces/" + url.PathEscape(namespace) +
		"/services/" + url.PathEscape(name) + ":" + port + "/proxy/" + restPath
	if encoded := forwardQuery.Encode(); encoded != "" {
		target += "?" + encoded
	}

	transport, err := rest.TransportFor(client.Config)
	if err != nil {
		c.String(http.StatusBadGateway, "cluster transport: %v", err)
		return
	}
	req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, target, c.Request.Body)
	if err != nil {
		c.String(http.StatusBadGateway, "proxy request: %v", err)
		return
	}
	req.Header = c.Request.Header.Clone()
	req.Header.Del("Authorization")
	req.Header.Del("Cookie")
	resp, err := transport.RoundTrip(req)
	if err != nil {
		c.String(http.StatusBadGateway, "service proxy: %v", err)
		return
	}
	defer resp.Body.Close()

	prefix := "/api/v1/namespaces/" + url.PathEscape(namespace) + "/services/" + url.PathEscape(name) + "/proxy"
	contentType := resp.Header.Get("Content-Type")
	rewrite := strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/css")
	if loc := resp.Header.Get("Location"); strings.HasPrefix(loc, "/") && !strings.HasPrefix(loc, "//") {
		resp.Header.Set("Location", prefix+loc)
	}

	var rewritten []byte
	if rewrite {
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			c.String(http.StatusBadGateway, "read upstream: %v", err)
			return
		}
		text := rewriteRootPaths(string(raw), prefix)
		if strings.Contains(contentType, "text/html") {
			text = injectBase(text, prefix)
		}
		rewritten = []byte(text)
		resp.Header.Del("Content-Length")
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	c.Status(resp.StatusCode)
	if rewrite {
		_, _ = c.Writer.Write(rewritten)
		return
	}
	_, _ = io.Copy(c.Writer, resp.Body)
}

// showcaseServiceHTML is a static page for the public exhibit. Root-relative
// links go through the same rewrite as a real Service response.
func showcaseServiceHTML(namespace, name, port string) string {
	return `<!doctype html>
<html>
<head><title>` + htmlEscape(name) + `</title></head>
<body style="font-family: sans-serif; margin: 2rem; line-height: 1.5">
  <p style="letter-spacing: .12em; font-size: 12px">SHOWCASE SERVICE</p>
  <h1>` + htmlEscape(name) + `</h1>
  <p>Namespace ` + htmlEscape(namespace) + `, port ` + htmlEscape(port) + `.</p>
  <p>This page is served by the exhibit. It is not a live workload.</p>
  <p><a href="/health">Health</a></p>
</body>
</html>`
}

func htmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return r.Replace(s)
}

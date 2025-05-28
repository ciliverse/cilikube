package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"

	"github.com/ciliverse/cilikube/internal/service"
)

type ProxyHandler struct {
	service *service.ProxyService
}

func NewProxyHandler(service *service.ProxyService) *ProxyHandler {
	return &ProxyHandler{service: service}
}

func (p *ProxyHandler) Proxy(c *gin.Context) {
	config := p.service.GetConfig()
	transport, err := rest.TransportFor(config)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "服务器内部错误: "+err.Error())
		return
	}
	target, err := p.validateTarget(*c.Request.URL, config.Host)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "服务器内部错误: "+err.Error())
		return
	}
	httpProxy := proxy.NewUpgradeAwareHandler(target, transport, false, false, nil)
	httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(transport, transport)
	httpProxy.ServeHTTP(c.Writer, c.Request)
}

func (p *ProxyHandler) validateTarget(target url.URL, host string) (*url.URL, error) {
	kubeURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	target.Path = target.Path[len("/api/v1/proxy/"):]

	target.Host = kubeURL.Host
	target.Scheme = kubeURL.Scheme
	return &target, nil
}

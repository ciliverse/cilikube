package handlers

import (
	"net/http"
	"net/url"

	"github.com/ciliverse/cilikube/pkg/k8s"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
)

type ProxyHandler struct {
	clusterManager *k8s.ClusterManager
}

func NewProxyHandler(cm *k8s.ClusterManager) *ProxyHandler {
	return &ProxyHandler{clusterManager: cm}
}

func (h *ProxyHandler) Proxy(c *gin.Context) {
	k8sClient, ok := k8s.GetK8sClientFromContext(c, h.clusterManager)
	if !ok {
		return
	}

	config := k8sClient.Config
	transport, err := rest.TransportFor(config)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "服务器内部错误: "+err.Error())
		return
	}
	target, err := h.validateTarget(*c.Request.URL, config.Host)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "服务器内部错误: "+err.Error())
		return
	}
	httpProxy := proxy.NewUpgradeAwareHandler(target, transport, false, false, nil)
	httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(transport, transport)
	httpProxy.ServeHTTP(c.Writer, c.Request)
}

func (h *ProxyHandler) validateTarget(target url.URL, host string) (*url.URL, error) {
	kubeURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	target.Path = target.Path[len("/api/v1/proxy/"):]

	target.Host = kubeURL.Host
	target.Scheme = kubeURL.Scheme
	return &target, nil
}

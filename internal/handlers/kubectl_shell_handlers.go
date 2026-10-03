package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/auth"
	"github.com/ciliverse/cilikube/pkg/k8s"
	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// KubectlShellHandler runs a local kubectl against the selected cluster.
// The process uses the cluster credential stored in CiliKube, so only admins may open it.
type KubectlShellHandler struct {
	clusterManager *k8s.ClusterManager
	audit          *service.AuditService
	upgrader       websocket.Upgrader
}

func NewKubectlShellHandler(cm *k8s.ClusterManager, audit *service.AuditService) *KubectlShellHandler {
	return &KubectlShellHandler{
		clusterManager: cm,
		audit:          audit,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (h *KubectlShellHandler) Shell(c *gin.Context) {
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("kubectl shell upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	if !adminRole(c) {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("admin role required\r\n"))
		return
	}
	client, ok := k8s.GetClientFromQuery(c, h.clusterManager)
	if !ok {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("failed to get Kubernetes client\r\n"))
		return
	}
	if k8s.IsShowcaseConfig(client.Config) {
		logShell(h.audit, c, "kubectl", "showcase")
		service.RunShowcaseTerminal(ws, service.ShowcaseKubectlBanner(), "demo:~$ ", service.ShowcaseKubectl)
		return
	}
	if _, err := exec.LookPath("kubectl"); err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("kubectl is not installed on the CiliKube server\r\n"))
		return
	}
	raw, err := k8s.KubeconfigFromREST(client.Config)
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("kubeconfig: %v\r\n", err)))
		return
	}
	file, err := os.CreateTemp("", "cilikube-kubeconfig-")
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("temp file: %v\r\n", err)))
		return
	}
	kubeconfigPath := file.Name()
	defer os.Remove(kubeconfigPath)
	if err := file.Chmod(0o600); err != nil {
		file.Close()
		_ = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("chmod: %v\r\n", err)))
		return
	}
	if _, err := file.Write(raw); err != nil {
		file.Close()
		_ = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("write kubeconfig: %v\r\n", err)))
		return
	}
	file.Close()

	cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 32, Cols: 120})
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("pty: %v\r\n", err)))
		return
	}
	defer func() {
		_ = ptmx.Close()
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()
	logShell(h.audit, c, "kubectl", "")

	go func() {
		buf := make([]byte, 4096)
		for {
			n, readErr := ptmx.Read(buf)
			if n > 0 {
				_ = ws.WriteMessage(websocket.BinaryMessage, append([]byte(nil), buf[:n]...))
			}
			if readErr != nil {
				return
			}
		}
	}()
	for {
		_, msg, readErr := ws.ReadMessage()
		if readErr != nil {
			return
		}
		if bytes.HasPrefix(bytes.TrimSpace(msg), []byte(`{"type":"resize"`)) {
			var resize struct {
				Type string `json:"type"`
				Cols uint16 `json:"cols"`
				Rows uint16 `json:"rows"`
			}
			if json.Unmarshal(msg, &resize) == nil && resize.Type == "resize" && resize.Cols > 0 && resize.Rows > 0 {
				_ = pty.Setsize(ptmx, &pty.Winsize{Cols: resize.Cols, Rows: resize.Rows})
				continue
			}
		}
		if _, writeErr := ptmx.Write(msg); writeErr != nil {
			return
		}
	}
}

func adminRole(c *gin.Context) bool {
	role, _ := c.Get("user_role")
	return role == "admin"
}

func logShell(audit *service.AuditService, c *gin.Context, action, target string) {
	if audit == nil {
		return
	}
	var userID uint
	if id, ok := c.Get("user_id"); ok {
		if n, ok := id.(uint); ok {
			userID = n
		}
	}
	username, _ := c.Get("username")
	name, _ := username.(string)
	_ = audit.LogResourceAccessEvent(userID, name, "shell", action, auth.AuditClientIP(c), c.GetHeader("User-Agent"), true, map[string]interface{}{
		"target": target,
	})
}

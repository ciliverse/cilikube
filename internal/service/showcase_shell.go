package service

import (
	"strings"
)

// ShowcaseKubectl answers a few read-only commands for the public exhibit.
// It never executes a local binary or contacts an apiserver.
func ShowcaseKubectl(line string) (string, bool) {
	cmd := normalizeShowcaseCmd(line)
	if cmd == "" {
		return "", false
	}
	if cmd == "exit" || cmd == "quit" {
		return "session closed\r\n", true
	}
	if blockedShowcaseCmd(cmd) {
		return "simulated shell: that command is not available on the public demo\r\n", false
	}
	switch cmd {
	case "help", "?", "kubectl", "kubectl help":
		return "" +
			"Simulated kubectl for the public demo. Read-only examples:\r\n" +
			"  kubectl version\r\n" +
			"  kubectl get nodes\r\n" +
			"  kubectl get pods -A\r\n" +
			"  kubectl get svc -n default\r\n" +
			"  exit\r\n", false
	case "kubectl version", "kubectl version --client", "version":
		return "Client Version: v1.36.2-showcase\r\nServer Version: v1.36.2-showcase\r\n", false
	case "kubectl get nodes", "kubectl get node", "kubectl get no":
		return "" +
			"NAME            STATUS   ROLES           VERSION\r\n" +
			"demo-master-1   Ready    control-plane   v1.36.2-showcase\r\n" +
			"demo-worker-1   Ready    worker          v1.36.2-showcase\r\n" +
			"demo-worker-2   Ready    worker          v1.36.2-showcase\r\n", false
	case "kubectl get pods", "kubectl get pod", "kubectl get po", "kubectl get pods -n default", "kubectl get pods --namespace default":
		return "" +
			"NAME                          READY   STATUS    NODE\r\n" +
			"web-frontend-7d9f8b-abc12     1/1     Running   demo-worker-1\r\n" +
			"web-frontend-7d9f8b-def34     1/1     Running   demo-worker-2\r\n" +
			"api-gateway-6c4d5-jkl78       1/1     Running   demo-worker-2\r\n", false
	case "kubectl get pods -A", "kubectl get pods --all-namespaces", "kubectl get pod -A":
		return "" +
			"NAMESPACE    NAME                          READY   STATUS    NODE\r\n" +
			"default      web-frontend-7d9f8b-abc12     1/1     Running   demo-worker-1\r\n" +
			"default      api-gateway-6c4d5-jkl78       1/1     Running   demo-worker-2\r\n" +
			"production   orders-api-8f2a1-aa111        1/1     Running   demo-worker-1\r\n" +
			"kube-system  coredns-xyz01                 1/1     Running   demo-master-1\r\n", false
	case "kubectl get svc", "kubectl get svc -n default", "kubectl get services", "kubectl get services -n default":
		return "" +
			"NAME           TYPE        CLUSTER-IP   PORT(S)\r\n" +
			"web-frontend   ClusterIP   10.0.12.20   80/TCP\r\n" +
			"api-gateway    ClusterIP   10.0.12.21   8080/TCP\r\n", false
	case "kubectl config current-context":
		return "demo\r\n", false
	default:
		return "simulated shell: that command is not available on the public demo\r\n", false
	}
}

// ShowcaseNodeShell answers a few read-only commands for one simulated node.
// It never starts a pod or reads the API host filesystem.
func ShowcaseNodeShell(node, line string) (string, bool) {
	node = strings.TrimSpace(node)
	if node == "" {
		node = "demo-worker-1"
	}
	cmd := normalizeShowcaseCmd(line)
	if cmd == "" {
		return "", false
	}
	if cmd == "exit" || cmd == "quit" {
		return "session closed\r\n", true
	}
	if blockedShowcaseCmd(cmd) {
		return "simulated shell: that command is not available on the public demo\r\n", false
	}
	switch cmd {
	case "help", "?":
		return "hostname, whoami, uname -a, pwd, exit\r\n", false
	case "hostname":
		return node + "\r\n", false
	case "whoami":
		return "root\r\n", false
	case "uname -a", "uname":
		return "Linux " + node + " 6.8.0-showcase x86_64 GNU/Linux\r\n", false
	case "pwd":
		return "/\r\n", false
	default:
		return "simulated shell: that command is not available on the public demo\r\n", false
	}
}

func ShowcaseKubectlBanner() string {
	return "Showcase cluster shell\r\n" +
		"This session is simulated. It does not run on the API host and cannot change a cluster.\r\n" +
		"Type help for the commands this demo accepts.\r\n"
}

func ShowcaseNodeBanner(node string) string {
	node = strings.TrimSpace(node)
	if node == "" {
		node = "demo-worker-1"
	}
	return "Showcase node shell on " + node + "\r\n" +
		"This session is simulated. No privileged pod is created.\r\n"
}

func ShowcaseNodePrompt(node string) string {
	node = strings.TrimSpace(node)
	if node == "" {
		node = "demo-worker-1"
	}
	return node + ":/# "
}

func normalizeShowcaseCmd(line string) string {
	fields := strings.Fields(strings.TrimSpace(line))
	return strings.Join(fields, " ")
}

func blockedShowcaseCmd(cmd string) bool {
	return strings.ContainsAny(cmd, ";|&$`<>\\") || strings.Contains(cmd, "$(")
}

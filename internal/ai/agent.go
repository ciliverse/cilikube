package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

type sendFn func(SSEEvent)

func Ready(cfg *configs.Config) bool {
	if cfg == nil || !cfg.AI.Enabled {
		return false
	}
	p := strings.ToLower(cfg.AI.Provider)
	if p == "" || p == "mock" {
		return true
	}
	return cfg.AI.APIKey != ""
}

func ProcessChat(ctx context.Context, cfg *configs.Config, client *k8s.Client, req *ChatRequest, send sendFn) {
	if cfg == nil || client == nil {
		send(SSEEvent{Event: "error", Data: map[string]string{"message": "AI or cluster unavailable"}})
		return
	}
	if len(req.Messages) == 0 {
		send(SSEEvent{Event: "error", Data: map[string]string{"message": "No messages"}})
		return
	}
	provider := strings.ToLower(cfg.AI.Provider)
	if provider == "" {
		provider = "mock"
	}
	switch provider {
	case "mock":
		processMock(ctx, client, req, send)
	case "openai":
		if err := processOpenAI(ctx, cfg, client, req, send); err != nil {
			send(SSEEvent{Event: "error", Data: map[string]string{"message": err.Error()}})
		}
	default:
		send(SSEEvent{Event: "error", Data: map[string]string{"message": fmt.Sprintf("unsupported provider %q", provider)}})
	}
}

func systemPrompt(namespace string) string {
	nsHint := "all namespaces"
	if namespace != "" {
		nsHint = "prefer namespace " + namespace
	}
	return strings.TrimSpace(`
You are CiliKube AI, a read-only Kubernetes assistant.
Use tools to gather evidence before answering. Never invent resource names.
Prefer concise Chinese answers unless the user writes in English.
When mentioning pods/deployments/services/nodes, the UI will show clickable cards — keep names accurate.
Scope hint: ` + nsHint + `
Only use these tools: get_cluster_overview, list_resources, get_resource, get_pod_logs.
Do not claim you can mutate the cluster.
`)
}

func lastUserText(msgs []ChatMessage) string {
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].Role == "user" {
			return msgs[i].Content
		}
	}
	return ""
}

func mergeRefs(dst *[]ResourceRef, add []ResourceRef) {
	seen := map[string]bool{}
	for _, r := range *dst {
		seen[r.Href] = true
	}
	for _, r := range add {
		if r.Href == "" || seen[r.Href] {
			continue
		}
		seen[r.Href] = true
		*dst = append(*dst, r)
	}
}

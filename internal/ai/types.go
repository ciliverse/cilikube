package ai

import "encoding/json"

// ChatRequest is the body for POST /api/v1/ai/chat.
type ChatRequest struct {
	Messages  []ChatMessage `json:"messages"`
	Namespace string        `json:"namespace,omitempty"`
	Language  string        `json:"language,omitempty"`
	// SkillID is optional telemetry from the UI skill chips (logging / mock routing).
	SkillID string `json:"skill_id,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ResourceRef is a clickable console deep-link emitted to the UI.
type ResourceRef struct {
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace,omitempty"`
	Name       string `json:"name"`
	Href       string `json:"href"`
	Console    string `json:"console,omitempty"` // logs | exec | attach
	Label      string `json:"label,omitempty"`
}

// SSEEvent is one server-sent event payload.
type SSEEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func MarshalSSE(ev SSEEvent) string {
	b, err := json.Marshal(ev.Data)
	if err != nil {
		b = []byte(`{"message":"encode error"}`)
	}
	return "event: " + ev.Event + "\ndata: " + string(b) + "\n\n"
}

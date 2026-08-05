package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

type oaiMsg struct {
	Role       string        `json:"role"`
	Content    string        `json:"content,omitempty"`
	ToolCalls  []oaiToolCall `json:"tool_calls,omitempty"`
	ToolCallID string        `json:"tool_call_id,omitempty"`
}

type oaiToolCall struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Function oaiFunction `json:"function"`
}

type oaiFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type oaiTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Parameters  map[string]interface{} `json:"parameters"`
	} `json:"function"`
}

func processOpenAI(ctx context.Context, cfg *configs.Config, client *k8s.Client, req *ChatRequest, send sendFn) error {
	if cfg.AI.APIKey == "" {
		return fmt.Errorf("AI API key not configured")
	}
	base := strings.TrimRight(cfg.AI.BaseURL, "/")
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	model := cfg.AI.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	lang := normalizeLang(req.Language)
	messages := []oaiMsg{{Role: "system", Content: systemPrompt(req.Namespace, lang)}}
	for _, m := range req.Messages {
		role := m.Role
		if role != "user" && role != "assistant" {
			role = "user"
		}
		messages = append(messages, oaiMsg{Role: role, Content: m.Content})
	}

	tools := make([]oaiTool, 0, len(toolDefinitions()))
	for _, t := range toolDefinitions() {
		var ot oaiTool
		ot.Type = "function"
		ot.Function.Name = t.Name
		ot.Function.Description = t.Description
		ot.Function.Parameters = t.Parameters
		tools = append(tools, ot)
	}

	var allRefs []ResourceRef
	httpClient := &http.Client{Timeout: 90 * time.Second}

	for i := 0; i < 6; i++ {
		body := map[string]interface{}{
			"model":       model,
			"messages":    messages,
			"tools":       tools,
			"tool_choice": "auto",
			"temperature": 0.2,
		}
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/chat/completions", bytes.NewReader(raw))
		if err != nil {
			return err
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+cfg.AI.APIKey)

		resp, err := httpClient.Do(httpReq)
		if err != nil {
			return err
		}
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("openai api %d: %s", resp.StatusCode, truncate(string(respBody), 400))
		}

		var parsed struct {
			Choices []struct {
				Message oaiMsg `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			return err
		}
		if len(parsed.Choices) == 0 {
			return fmt.Errorf("empty model response")
		}
		msg := parsed.Choices[0].Message
		if len(msg.ToolCalls) == 0 {
			content := strings.TrimSpace(msg.Content)
			if content == "" {
				if isZh(lang) {
					content = "（模型未返回文本）"
				} else {
					content = "(model returned no text)"
				}
			}
			if len(allRefs) > 0 {
				send(SSEEvent{Event: "resources", Data: map[string]interface{}{"items": allRefs}})
			}
			send(SSEEvent{Event: "message", Data: map[string]string{"content": content}})
			send(SSEEvent{Event: "done", Data: map[string]string{}})
			return nil
		}

		messages = append(messages, msg)
		for _, tc := range msg.ToolCalls {
			args := map[string]interface{}{}
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
			send(SSEEvent{Event: "tool_call", Data: map[string]interface{}{
				"name": tc.Function.Name, "arguments": args, "id": tc.ID,
			}})
			res, err := executeTool(ctx, client, tc.Function.Name, args, lang)
			content := ""
			ok := true
			if err != nil {
				ok = false
				content = err.Error()
			} else {
				content = res.Text
				mergeRefs(&allRefs, res.Resources)
			}
			send(SSEEvent{Event: "tool_result", Data: map[string]interface{}{
				"name": tc.Function.Name, "ok": ok, "content": truncate(content, 4000), "id": tc.ID,
			}})
			messages = append(messages, oaiMsg{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    content,
			})
		}
	}
	return fmt.Errorf("tool loop exceeded max iterations")
}

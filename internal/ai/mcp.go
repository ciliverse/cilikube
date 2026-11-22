package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

const (
	mcpProtocolVersion = "2025-03-26"
	releaseVersion     = "1.7.0"
)

// MCP serves a stateless, read-only Model Context Protocol endpoint.
// POST is JSON-RPC (initialize, tools/list, tools/call). GET describes the endpoint.
func (h *Handler) MCP(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.JSON(http.StatusOK, gin.H{
			"name":        "cilikube",
			"version":     releaseVersion,
			"protocol":    mcpProtocolVersion,
			"transport":   "streamable-http",
			"read_only":   true,
			"method":      "POST",
			"cluster_arg": "clusterId",
			"auth":        "Authorization: Bearer <token>",
		})
		return
	}

	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		writeRPCError(c, nil, -32700, "failed to read body")
		return
	}
	if len(body) > 0 && body[0] == '[' {
		writeRPCError(c, nil, -32600, "JSON-RPC batches are not supported")
		return
	}

	var req rpcRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeRPCError(c, nil, -32700, "invalid JSON")
		return
	}
	if req.JSONRPC != "" && req.JSONRPC != "2.0" {
		writeRPCError(c, req.ID, -32600, "jsonrpc must be 2.0")
		return
	}
	if isRPCNotification(req.ID) || strings.HasPrefix(req.Method, "notifications/") {
		c.Status(http.StatusAccepted)
		return
	}

	lang := c.Query("lang")
	if lang == "" {
		lang = c.GetHeader("Accept-Language")
	}

	var client *k8s.Client
	if req.Method == "tools/call" {
		var ok bool
		client, ok = k8s.GetClientFromQuery(c, h.manager)
		if !ok {
			return
		}
	}

	result, rpcErr := dispatchMCP(c.Request.Context(), req, client, h.prom, lang)
	if rpcErr != nil {
		writeRPCError(c, req.ID, rpcErr.Code, rpcErr.Message)
		return
	}
	if sid := c.GetHeader("Mcp-Session-Id"); sid != "" {
		c.Header("Mcp-Session-Id", sid)
	}
	c.JSON(http.StatusOK, rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: result})
}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcErrorBody   `json:"error,omitempty"`
}

type rpcErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func isRPCNotification(id json.RawMessage) bool {
	return len(id) == 0 || string(id) == "null"
}

func writeRPCError(c *gin.Context, id json.RawMessage, code int, message string) {
	if len(id) == 0 {
		id = json.RawMessage("null")
	}
	c.JSON(http.StatusOK, rpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &rpcErrorBody{Code: code, Message: message},
	})
}

func dispatchMCP(ctx context.Context, req rpcRequest, client *k8s.Client, prom *service.PrometheusService, lang string) (any, *rpcErrorBody) {
	switch req.Method {
	case "initialize":
		return map[string]any{
			"protocolVersion": mcpProtocolVersion,
			"capabilities": map[string]any{
				"tools": map[string]any{"listChanged": false},
			},
			"serverInfo": map[string]any{
				"name":    "cilikube",
				"version": releaseVersion,
			},
			"instructions": "CiliKube investigator tools are read-only. They cannot create, update, delete, scale, or exec. clusterId selects the cluster; the Bearer token is the logged-in user.",
		}, nil
	case "ping", "logging/setLevel":
		return map[string]any{}, nil
	case "tools/list":
		tools := make([]map[string]any, 0, len(toolDefinitions()))
		for _, t := range toolDefinitions() {
			tools = append(tools, map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"inputSchema": t.Parameters,
			})
		}
		return map[string]any{"tools": tools}, nil
	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if len(req.Params) > 0 {
			if err := json.Unmarshal(req.Params, &params); err != nil {
				return nil, &rpcErrorBody{Code: -32602, Message: "invalid tools/call params"}
			}
		}
		if params.Name == "" {
			return nil, &rpcErrorBody{Code: -32602, Message: "tool name is required"}
		}
		res, err := executeTool(ctx, client, prom, params.Name, params.Arguments, lang)
		text := ""
		isError := false
		if err != nil {
			isError = true
			text = err.Error()
		} else {
			text = res.Text
		}
		return map[string]any{
			"content": []map[string]any{{"type": "text", "text": text}},
			"isError": isError,
		}, nil
	default:
		return nil, &rpcErrorBody{Code: -32601, Message: "method not found"}
	}
}

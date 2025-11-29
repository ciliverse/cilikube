package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

type InstallerHandler struct {
	installerService service.InstallerService
}

func NewInstallerHandler(is service.InstallerService) *InstallerHandler {
	return &InstallerHandler{
		installerService: is,
	}
}

// HealthCheck handles health check requests
func (h *InstallerHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Backend service is running",
	})
}

// StreamMinikubeInstallation handles the SSE request.
func (h *InstallerHandler) StreamMinikubeInstallation(c *gin.Context) {
	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	// CORS handled by middleware

	// Flush headers
	c.Writer.Flush()

	// Create channel
	messageChan := make(chan service.ProgressUpdate) // Use service.ProgressUpdate

	// Get client disconnect notification channel
	// Using c.Request.Context().Done() is more modern and recommended
	clientGone := c.Request.Context().Done() // Type is <-chan struct{}

	log.Println("SSE: Connection established, starting installation service Goroutine.")
	// Start service in new Goroutine
	go h.installerService.InstallMinikube(messageChan, clientGone)
	// Pass clientGone (<-chan struct{}) to service

	log.Println("SSE: Handler starts listening to service messages and pushing to client...")
	// Process stream in current Goroutine until completion or error
	err := h.streamUpdatesToClient(c, messageChan, clientGone) // Pass clientGone (<-chan struct{}) to helper function
	if err != nil {
		log.Printf("SSE: Stream processing error: %v", err)
	}
	log.Println("SSE: Handler stream processing ended.")
}

// streamUpdatesToClient helper function that processes messages from service and pushes to client
func (h *InstallerHandler) streamUpdatesToClient(c *gin.Context, messageChan <-chan service.ProgressUpdate, clientGone <-chan struct{}) error {
	defer log.Println("SSE: streamUpdatesToClient loop ended.")
	for {
		select {
		case <-clientGone: // Listen to Context.Done() channel
			log.Println("SSE: Client disconnected (Context Done).")
			return nil // Client disconnected, normal exit
		case update, ok := <-messageChan:
			if !ok {
				log.Println("SSE: Service channel closed.")
				return nil // Service completed or error, normal exit from loop
			}

			// Received update, prepare to send
			log.Printf("SSE: Received update from service: Step=%s, Progress=%d, Done=%t", update.Step, update.Progress, update.Done)

			jsonData, err := json.Marshal(update)
			if err != nil {
				log.Printf("SSE: Failed to serialize service update: %v", err)
				// Try to notify client
				_, writeErr := fmt.Fprintf(c.Writer, "event: error\ndata: {\"error\": \"Internal server error marshalling update: %v\"}\n\n", err)
				if writeErr != nil {
					log.Printf("SSE: Failed to write serialization error to client: %v", writeErr)
					return writeErr // Return write error
				}
				c.Writer.Flush()
				continue // Continue listening for next message
			}

			// Send data
			_, writeErr := fmt.Fprintf(c.Writer, "event: message\ndata: %s\n\n", string(jsonData))
			if writeErr != nil {
				log.Printf("SSE: Failed to write data to client: %v", writeErr)
				return writeErr // Return write error
			}

			// Flush to ensure sending
			if f, ok := c.Writer.(http.Flusher); ok {
				f.Flush()
			} else {
				log.Println("SSE: Warning - ResponseWriter does not support Flushing.")
			}

			// Exit if this is the last message
			if update.Done {
				log.Println("SSE: Final update sent, normal exit from stream processing.")
				return nil
			}
		}
	}
}

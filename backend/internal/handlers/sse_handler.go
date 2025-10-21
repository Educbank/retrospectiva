package handlers

import (
	"net/http"
	"time"

	"educ-retro/internal/auth"
	"educ-retro/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SSEHandler struct {
	realtimeService *services.RealtimeService
}

func NewSSEHandler(realtimeService *services.RealtimeService) *SSEHandler {
	return &SSEHandler{
		realtimeService: realtimeService,
	}
}

func (h *SSEHandler) HandleSSE(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	// Validate token
	claims, err := auth.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	retrospectiveIDStr := c.Query("retrospective_id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// Register client
	client := h.realtimeService.RegisterClient(retrospectiveID)
	defer h.realtimeService.UnregisterClient(client)

	// Send initial connection message
	c.SSEvent("message", map[string]interface{}{
		"type": "connected",
		"data": map[string]interface{}{
			"user_id":          claims.UserID,
			"user_name":        claims.Name,
			"retrospective_id": retrospectiveID,
		},
		"timestamp": time.Now().Unix(),
	})
	c.Writer.Flush()

	// Send current blur state to new client
	blurState := h.realtimeService.GetBlurState(retrospectiveID)
	if blurState {
		c.SSEvent("message", map[string]interface{}{
			"type": "blur_toggled",
			"data": map[string]interface{}{
				"blurred": blurState,
			},
			"timestamp": time.Now().Unix(),
		})
		c.Writer.Flush()
	}

	// Keep connection alive and send events
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			// Client disconnected
			return
		case <-ticker.C:
			// Send keepalive
			c.SSEvent("ping", map[string]interface{}{
				"timestamp": time.Now().Unix(),
			})
			c.Writer.Flush()
		default:
			// Try to get event from client channel
			select {
			case eventData := <-client.Send:
				c.SSEvent("message", eventData)
				c.Writer.Flush()
			case <-time.After(1 * time.Second):
				// No event, continue
			}
		}
	}
}

func (h *SSEHandler) SetupRoutes(r *gin.RouterGroup) {
	sse := r.Group("/sse")
	{
		sse.GET("/retrospective", h.HandleSSE)
	}
}

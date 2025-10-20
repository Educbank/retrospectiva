package handlers

import (
	"log"
	"net/http"

	"educ-retro/internal/auth"
	ws "educ-retro/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		return true
	},
}

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	// Validate token manually
	claims, err := auth.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userUUID := claims.UserID
	userName := claims.Name

	retrospectiveIDStr := c.Query("retrospective_id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create client
	client := &ws.Client{
		ID:              uuid.New(),
		UserID:          userUUID,
		UserName:        userName,
		RetrospectiveID: retrospectiveID,
		Conn:            conn,
		Send:            make(chan []byte, 256),
		Hub:             h.hub,
	}

	// Register client with hub
	h.hub.Register <- client

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// GetRetrospectiveParticipants returns the list of participants in a retrospective
func (h *WebSocketHandler) GetRetrospectiveParticipants(c *gin.Context) {
	retrospectiveIDStr := c.Param("id")
	retrospectiveID, err := uuid.Parse(retrospectiveIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retrospective ID"})
		return
	}

	clients := h.hub.GetRoomClients(retrospectiveID)

	var participants []map[string]interface{}
	for _, client := range clients {
		participants = append(participants, map[string]interface{}{
			"user_id": client.UserID,
			"name":    client.UserName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"participants": participants,
		"count":        len(participants),
	})
}

func (h *WebSocketHandler) SetupRoutes(r *gin.RouterGroup) {
	ws := r.Group("/ws")
	{
		ws.GET("/retrospective", h.HandleWebSocket)
		ws.GET("/retrospective/:id/participants", h.GetRetrospectiveParticipants)
	}
}

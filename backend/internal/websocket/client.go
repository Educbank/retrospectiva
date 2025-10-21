package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"educ-retro/internal/models"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, you should check the origin
		return true
	},
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var message models.WebSocketMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		c.handleMessage(message)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(message models.WebSocketMessage) {
	// Set the user ID from the authenticated client
	message.UserID = &c.UserID

	switch message.Type {
	case models.WSMessageJoinRetrospective:
		// Already handled in the connection setup
		
	case models.WSMessageLeaveRetrospective:
		c.Hub.Unregister <- c
		
	case models.WSMessageNewItem:
		// This would typically involve saving to database
		// For now, just broadcast to the room
		c.Hub.BroadcastToRetrospective(c.RetrospectiveID, message)
		
	case models.WSMessageUpdateItem:
		c.Hub.BroadcastToRetrospective(c.RetrospectiveID, message)
		
	case models.WSMessageDeleteItem:
		c.Hub.BroadcastToRetrospective(c.RetrospectiveID, message)
		
	case models.WSMessageVoteItem:
		c.Hub.BroadcastToRetrospective(c.RetrospectiveID, message)
		
	case models.WSMessageUnvoteItem:
		c.Hub.BroadcastToRetrospective(c.RetrospectiveID, message)
		
	case models.WSMessageUpdateRetrospective:
		c.Hub.BroadcastToRetrospective(c.RetrospectiveID, message)
		
	case models.WSMessageNewActionItem:
		c.Hub.BroadcastToRetrospective(c.RetrospectiveID, message)
		
	case models.WSMessageUpdateActionItem:
		c.Hub.BroadcastToRetrospective(c.RetrospectiveID, message)
		
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

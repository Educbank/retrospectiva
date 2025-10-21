package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"educ-retro/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	UserName       string
	RetrospectiveID uuid.UUID
	Conn           *websocket.Conn
	Send           chan []byte
	Hub            *Hub
}

type Hub struct {
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	broadcast  chan []byte
	rooms      map[uuid.UUID][]*Client // retrospective_id -> clients
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[uuid.UUID][]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	
	// Add client to retrospective room
	h.rooms[client.RetrospectiveID] = append(h.rooms[client.RetrospectiveID], client)

	// Notify other clients in the same retrospective
	userJoinedMsg := models.WebSocketMessage{
		Type: models.WSMessageUserJoined,
		Data: json.RawMessage(`{"user_id":"` + client.UserID.String() + `","name":"` + client.UserName + `"}`),
		Timestamp: time.Now().Unix(),
	}

	h.broadcastToRoom(client.RetrospectiveID, userJoinedMsg)

	log.Printf("Client %s joined retrospective %s", client.UserName, client.RetrospectiveID)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		// Remove client from retrospective room
		room := h.rooms[client.RetrospectiveID]
		for i, c := range room {
			if c == client {
				h.rooms[client.RetrospectiveID] = append(room[:i], room[i+1:]...)
				break
			}
		}

		// If room is empty, delete it
		if len(h.rooms[client.RetrospectiveID]) == 0 {
			delete(h.rooms, client.RetrospectiveID)
		}

		// Notify other clients in the same retrospective
		userLeftMsg := models.WebSocketMessage{
			Type: models.WSMessageUserLeft,
			Data: json.RawMessage(`{"user_id":"` + client.UserID.String() + `"}`),
			Timestamp: time.Now().Unix(),
		}

		h.broadcastToRoom(client.RetrospectiveID, userLeftMsg)

		log.Printf("Client %s left retrospective %s", client.UserName, client.RetrospectiveID)
	}
}

func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) BroadcastToRetrospective(retrospectiveID uuid.UUID, message models.WebSocketMessage) {
	h.broadcastToRoom(retrospectiveID, message)
}

func (h *Hub) broadcastToRoom(retrospectiveID uuid.UUID, message models.WebSocketMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	clients := h.rooms[retrospectiveID]
	for _, client := range clients {
		select {
		case client.Send <- messageBytes:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) GetRoomClients(retrospectiveID uuid.UUID) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[retrospectiveID]
}

func (h *Hub) GetRoomClientCount(retrospectiveID uuid.UUID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[retrospectiveID])
}

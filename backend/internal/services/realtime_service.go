package services

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RealtimeEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type RealtimeService struct {
	clients    map[string]chan RealtimeEvent
	clientsMu  sync.RWMutex
	register   chan *RealtimeClient
	unregister chan *RealtimeClient
	broadcast  chan RealtimeEvent
	blurStates map[uuid.UUID]bool // Map to store blur state per retrospective
	blurMu     sync.RWMutex
}

type RealtimeClient struct {
	ID              string
	RetrospectiveID uuid.UUID
	Send            chan RealtimeEvent
	Service         *RealtimeService
}

func NewRealtimeService() *RealtimeService {
	service := &RealtimeService{
		clients:    make(map[string]chan RealtimeEvent),
		register:   make(chan *RealtimeClient),
		unregister: make(chan *RealtimeClient),
		broadcast:  make(chan RealtimeEvent),
		blurStates: make(map[uuid.UUID]bool),
	}

	go service.run()
	return service
}

func (s *RealtimeService) run() {
	for {
		select {
		case client := <-s.register:
			s.clientsMu.Lock()
			s.clients[client.ID] = client.Send
			s.clientsMu.Unlock()

		case client := <-s.unregister:
			s.clientsMu.Lock()
			if send, ok := s.clients[client.ID]; ok {
				close(send)
				delete(s.clients, client.ID)
			}
			s.clientsMu.Unlock()

		case event := <-s.broadcast:
			s.clientsMu.RLock()
			for _, send := range s.clients {
				select {
				case send <- event:
				default:
					// Client is not ready, skip
				}
			}
			s.clientsMu.RUnlock()
		}
	}
}

func (s *RealtimeService) RegisterClient(retrospectiveID uuid.UUID) *RealtimeClient {
	clientID := uuid.New().String()
	send := make(chan RealtimeEvent, 256)

	client := &RealtimeClient{
		ID:              clientID,
		RetrospectiveID: retrospectiveID,
		Send:            send,
		Service:         s,
	}

	s.register <- client
	return client
}

func (s *RealtimeService) UnregisterClient(client *RealtimeClient) {
	s.unregister <- client
}

func (s *RealtimeService) BroadcastToRetrospective(retrospectiveID uuid.UUID, eventType string, data interface{}) {
	event := RealtimeEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	s.clientsMu.RLock()
	for clientID, send := range s.clients {
		// Get client from service to check retrospective ID
		// For now, broadcast to all clients (we can optimize later)
		select {
		case send <- event:
		default:
			// Remove client if channel is full
			delete(s.clients, clientID)
			close(send)
		}
	}
	s.clientsMu.RUnlock()
}

func (c *RealtimeClient) SendJSON() ([]byte, error) {
	select {
	case event := <-c.Send:
		return json.Marshal(event)
	case <-time.After(30 * time.Second):
		// Send keepalive
		return json.Marshal(RealtimeEvent{
			Type:      "ping",
			Data:      nil,
			Timestamp: time.Now().Unix(),
		})
	}
}

// SetBlurState sets the blur state for a retrospective
func (s *RealtimeService) SetBlurState(retrospectiveID uuid.UUID, blurred bool) {
	s.blurMu.Lock()
	defer s.blurMu.Unlock()
	s.blurStates[retrospectiveID] = blurred
}

// GetBlurState gets the blur state for a retrospective
func (s *RealtimeService) GetBlurState(retrospectiveID uuid.UUID) bool {
	s.blurMu.RLock()
	defer s.blurMu.RUnlock()
	return s.blurStates[retrospectiveID]
}

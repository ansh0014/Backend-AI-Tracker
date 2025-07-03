package ws

import (
	"Tracker/internal/model"
	"sync"
)

type Manager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.mu.Unlock()
		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
			m.mu.Unlock()
		case message := <-m.broadcast:
			m.mu.Lock()
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, client)
				}
			}
			m.mu.Unlock()
		}
	}
}

func (m *Manager) RegisterClient(client *Client) {
	m.register <- client
	client.manager = m
}

func (m *Manager) UnregisterClient(client *Client) {
	m.unregister <- client
}

func (m *Manager) Broadcast(event *model.Event) {
	// Marshal event to JSON and broadcast (implementation can be added)
}

func (m *Manager) ProcessEvent(event *model.Event) {
	// Placeholder for event processing logic (see processor.go)
}

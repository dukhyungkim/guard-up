package events

import (
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

var eventManagerInstance *EventManager

type Event struct {
	Type EventType
	Data string
}

type EventManager struct {
	Message      chan Event
	NewClient    chan chan Event
	CloseClient  chan chan Event
	TotalClients map[chan Event]struct{}
}

func NewEventManager() *EventManager {
	if eventManagerInstance != nil {
		return eventManagerInstance
	}

	eventManagerInstance = &EventManager{
		Message:      make(chan Event),
		NewClient:    make(chan chan Event),
		CloseClient:  make(chan chan Event),
		TotalClients: make(map[chan Event]struct{}),
	}
	return eventManagerInstance
}

func (m *EventManager) HandleEvent() {
	for {
		select {
		case client := <-m.NewClient:
			m.TotalClients[client] = struct{}{}
			log.Printf("new client connected. total: %d\n", len(m.TotalClients))

		case client := <-m.CloseClient:
			delete(m.TotalClients, client)
			log.Printf("client closed. total: %d\n", len(m.TotalClients))

		case event := <-m.Message:
			for client := range m.TotalClients {
				client <- event
			}
		}
	}
}

func (m *EventManager) HandleSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")

	stream := make(chan Event)
	m.NewClient <- stream
	defer func() {
		m.CloseClient <- stream
		close(stream)
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-stream; ok {
			c.SSEvent(msg.Type.String(), msg.Data)
			return true
		}
		return false
	})
}

func (m *EventManager) SendEvent(eventType EventType, eventData any) {
	data, _ := json.Marshal(eventData)
	m.Message <- Event{
		Type: eventType,
		Data: string(data),
	}
}

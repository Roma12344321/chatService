package ws

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[int]map[*Client]bool // Клиенты сгруппированы по RoomId
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	Id       int    `json:"id"`
	PersonId int    `json:"person_id"`
	RoomId   int    `json:"room_id"`
	Content  string `json:"content"`
	Method   string `json:"method"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.roomId] == nil {
				h.clients[client.roomId] = make(map[*Client]bool)
			}
			h.clients[client.roomId][client] = true
		case client := <-h.unregister:
			if clients, ok := h.clients[client.roomId]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.roomId)
					}
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients[message.RoomId] {
				m, err := json.Marshal(message)
				if err != nil {
					log.Printf("error marshaling message: %v", err)
					continue
				}
				select {
				case client.send <- m:
				default:
					close(client.send)
					delete(h.clients[message.RoomId], client)
					if len(h.clients[message.RoomId]) == 0 {
						delete(h.clients, message.RoomId)
					}
				}
			}
		}
	}
}

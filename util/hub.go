package util

import (
	"encoding/json"
	"log"
)

type Hub struct {
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Rooms      map[string]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			room := h.Rooms[client.room]
			if room == nil {
				room = make(map[*Client]bool)
				h.Rooms[client.room] = room
			}
			room[client] = true
		case client := <-h.Unregister:
			room := h.Rooms[client.room]
			if room != nil {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.send)
				}
				if len(room) == 0 {
					delete(h.Rooms, client.room)
				}
			}
		case message := <-h.Broadcast:
			room := h.Rooms[message.Room]
			messageData, err := json.Marshal(message)
			if err != nil {
				log.Fatal(err)
				break
			}
			messageStr := string(messageData)

			for client := range room {
				select {
				case client.send <- []byte(messageStr):
				default:
					close(client.send)
					delete(room, client)
				}
			}
			if len(room) == 0 {
				delete(h.Rooms, message.Room)
			}
		}
	}
}

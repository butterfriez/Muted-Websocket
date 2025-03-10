package util

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 2048
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	room string
	send chan []byte // in JSON
}

type Message struct {
	Data string `json:"data"`
	Room string `json:"room"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}
		log.Printf(msg.Data)
		c.hub.Broadcast <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Write the first message
			w.Write([]byte(message))

			// Write queued chat messages in the same frame
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n")) // Newline separator
				w.Write([]byte(<-c.send))
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("x-auth-username")
	serverId := r.Header.Get("x-auth-serverId")
	room := r.Header.Get("x-auth-room")
	sessionToken := r.Header.Get("x-auth-sessionToken")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	authenticated, newToken := VerifyUser(username, serverId, sessionToken, conn)
	if !authenticated {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
	}

	w.Header().Add("x-auth-sessionToken", newToken)

	client := &Client{hub: hub, conn: conn, room: room, send: make(chan []byte, 256)}
	client.hub.Register <- client

	go client.writePump()
	go client.readPump()
	log.Print("servews new connection")
}

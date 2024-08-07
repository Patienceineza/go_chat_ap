package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type WebSocketServer struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Message struct {
	Type     string    `json:"type"`
	Content  string    `json:"content"`
	RoomID   uuid.UUID `json:"room_id,omitempty"`
	SenderID uuid.UUID `json:"sender_id"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Conn: conn, Send: make(chan []byte)}
	server.Register <- client

	go client.readPump(server)
	go client.writePump()
}

func (c *Client) readPump(server *WebSocketServer) {
	defer func() {
		server.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println(err)
			continue
		}

		server.Broadcast <- message
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (server *WebSocketServer) HandleMessages() {
	for {
		message := <-server.Broadcast
		for client := range server.Clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(server.Clients, client)
			}
		}
	}
}

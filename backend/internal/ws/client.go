package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer (increased for code updates)
	maxMessageSize = 65536
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   int64
	roomCode string
}

// IncomingMessage represents a message from client
type IncomingMessage struct {
	Type     string          `json:"type"`
	RoomCode string          `json:"room_code,omitempty"`
	UserID   int64           `json:"user_id,omitempty"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}

// CodeUpdate represents a code change from a player
type CodeUpdate struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Code     string `json:"code"`
	Cursor   int    `json:"cursor"`
}

// PlayerStatus represents player's current activity
type PlayerStatus struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Status   string `json:"status"` // typing, idle, submitting
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws: read error: %v", err)
			}
			break
		}

		// Parse incoming message
		var msg IncomingMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			continue
		}

		// Handle different message types
		switch msg.Type {
		case "join_room":
			c.userID = msg.UserID
			c.roomCode = msg.RoomCode
			c.hub.JoinRoom(c, msg.RoomCode)
			log.Printf("ws: [JOIN] client %d joined room %s (room has %d clients)", msg.UserID, msg.RoomCode, c.hub.RoomClientCount(msg.RoomCode))

		case "code_update":
			// Broadcast code update to room
			if c.roomCode != "" {
				log.Printf("ws: [CODE_UPDATE] from user %d in room %s, payload size: %d bytes", c.userID, c.roomCode, len(msg.Payload))
				if err := c.hub.BroadcastToRoom(c.roomCode, TypeCodeUpdate, msg.Payload); err != nil {
					log.Printf("ws: [ERROR] failed to broadcast code_update: %v", err)
				} else {
					log.Printf("ws: [CODE_UPDATE] broadcasted to %d clients in room %s", c.hub.RoomClientCount(c.roomCode), c.roomCode)
				}
			} else {
				log.Printf("ws: [WARN] code_update from user %d but no room set", c.userID)
			}

		case "player_status":
			// Broadcast player status to room
			if c.roomCode != "" {
				log.Printf("ws: [PLAYER_STATUS] from user %d in room %s", c.userID, c.roomCode)
				c.hub.BroadcastToRoom(c.roomCode, TypePlayerStatus, msg.Payload)
			}

		default:
			log.Printf("ws: [UNKNOWN] message type '%s' from user %d", msg.Type, c.userID)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
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

// ServeWs handles websocket requests from the peer
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws: upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

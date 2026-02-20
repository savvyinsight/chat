package ws

import (
	"log"
	"time"

	"chat/model"
	"chat/service"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Message represents a chat message.
type Message struct {
	Type   string `json:"type"`
	From   uint   `json:"from"`
	To     uint   `json:"to,omitempty"`
	RoomID string `json:"room_id,omitempty"`
	ID     uint   `json:"id,omitempty"`
	Body   string `json:"body"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *Message

	// user id associated with this connection
	userID uint
}

func NewClient(h *Hub, conn *websocket.Conn, userID uint) *Client {
	return &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan *Message, 256),
		userID: userID,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		log.Printf("pong received from user %d", c.userID)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var msg Message
		if err := c.conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// set sender
		msg.From = c.userID

		// handle join/leave room messages
		if msg.Type == "join" && msg.RoomID != "" {
			c.hub.joinRoom(msg.RoomID, c)
			continue
		}
		if msg.Type == "leave" && msg.RoomID != "" {
			c.hub.leaveRoom(msg.RoomID, c)
			continue
		}

		// handle ack messages
		if msg.Type == "ack" && msg.ID != 0 {
			// mark message delivered
			if err := service.AckMessage(msg.ID); err != nil {
				log.Printf("ack update failed: %v", err)
			}
			continue
		}

		// persist message to DB
		mm := &model.Message{
			From: msg.From,
			To:   msg.To,
			Room: msg.RoomID,
			Type: msg.Type,
			Body: msg.Body,
		}
		if err := service.SaveMessage(mm); err != nil {
			log.Printf("save message failed: %v", err)
		} else {
			// set generated ID so receivers can ack
			msg.ID = mm.ID
		}

		c.hub.broadcast <- &msg
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ping failed for user %d: %v", c.userID, err)
				return
			}
			log.Printf("ping sent to user %d", c.userID)
		}
	}
}

package ws

import (
	"chat/global"
	"chat/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Map userID -> set of clients
	users map[uint]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Redis pubsub subscriptions per channel
	subs   map[string]*RedisSub
	subsMu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client, 128),
		unregister: make(chan *Client, 128),
		clients:    make(map[*Client]bool),
		users:      make(map[uint]map[*Client]bool),
		subs:       make(map[string]*RedisSub),
	}
}

// DefaultHub is the package-level hub used by the server.
var DefaultHub = NewHub()

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			if c.userID != 0 {
				if _, ok := h.users[c.userID]; !ok {
					h.users[c.userID] = make(map[*Client]bool)
				}
				h.users[c.userID][c] = true
				// ensure redis subscription for user channel
				h.ensureSub(fmt.Sprintf("user:%d", c.userID))
			}
			log.Printf("client registered: user=%d total=%d", c.userID, len(h.clients))
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			if c.userID != 0 {
				if set, ok := h.users[c.userID]; ok {
					delete(set, c)
					if len(set) == 0 {
						delete(h.users, c.userID)
						// remove redis subscription when no local clients
						h.removeSub(fmt.Sprintf("user:%d", c.userID))
					}
				}
			}
			log.Printf("client unregistered: user=%d total=%d", c.userID, len(h.clients))
		case m := <-h.broadcast:
			// publish to redis so other instances receive
			go h.publishToRedis(m)
			if m.To != 0 {
				// targeted to a user
				deliveredLocal := false
				if set, ok := h.users[m.To]; ok {
					for c := range set {
						select {
						case c.send <- m:
							deliveredLocal = true
						default:
							close(c.send)
							delete(h.clients, c)
						}
					}
				}
				if deliveredLocal {
					// mark message delivered and notify sender
					if m.ID != 0 {
						go func(id uint) {
							if err := service.AckMessage(id); err != nil {
								log.Printf("AckMessage error: %v", err)
							}
						}(m.ID)
					}
					// send ack to sender local clients
					ack := &Message{Type: "ack", ID: m.ID}
					if set, ok := h.users[m.From]; ok {
						for c := range set {
							select {
							case c.send <- ack:
							default:
								close(c.send)
								delete(h.clients, c)
							}
						}
					}
					// publish ack to sender channel so other instances can update
					if m.From != 0 && global.GVA_REDIS != nil {
						b, _ := json.Marshal(ack)
						if err := global.GVA_REDIS.Publish(context.Background(), fmt.Sprintf("user:%d", m.From), string(b)).Err(); err != nil {
							log.Printf("redis publish ack error: %v", err)
						}
					}
				}
			} else {
				// broadcast to all
				for c := range h.clients {
					select {
					case c.send <- m:
					default:
						close(c.send)
						delete(h.clients, c)
					}
				}
			}
		}
	}
}

// RedisSub holds pubsub and cancel function
type RedisSub struct {
	channel string
	cancel  context.CancelFunc
}

func (h *Hub) ensureSub(channel string) {
	if global.GVA_REDIS == nil {
		return
	}
	h.subsMu.Lock()
	defer h.subsMu.Unlock()
	if _, ok := h.subs[channel]; ok {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	pubsub := global.GVA_REDIS.Subscribe(ctx, channel)
	h.subs[channel] = &RedisSub{channel: channel, cancel: cancel}
	ch := pubsub.Channel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = pubsub.Close()
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				var m Message
				if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
					log.Printf("redis unmarshal error: %v", err)
					continue
				}
				// deliver to local clients for target user
				if strings.HasPrefix(msg.Channel, "user:") && m.Type == "ack" {
					// channel owner id
					parts := strings.Split(msg.Channel, ":")
					if len(parts) == 2 {
						if uid, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
							owner := uint(uid)
							// mark ack in DB
							if m.ID != 0 {
								if err := service.AckMessage(m.ID); err != nil {
									log.Printf("AckMessage error from redis ack: %v", err)
								}
							}
							// deliver ack to local owner clients
							if set, ok := h.users[owner]; ok {
								for c := range set {
									select {
									case c.send <- &m:
									default:
										close(c.send)
										delete(h.clients, c)
									}
								}
							}
						}
					}
				} else {
					// deliver to local clients for target user
					if m.To != 0 {
						if set, ok := h.users[m.To]; ok {
							deliveredLocal := false
							for c := range set {
								select {
								case c.send <- &m:
									deliveredLocal = true
								default:
									close(c.send)
									delete(h.clients, c)
								}
							}
							if deliveredLocal {
								// notify sender via redis ack so origin marks delivered
								if m.ID != 0 && m.From != 0 && global.GVA_REDIS != nil {
									ack := &Message{Type: "ack", ID: m.ID}
									b, _ := json.Marshal(ack)
									if err := global.GVA_REDIS.Publish(context.Background(), fmt.Sprintf("user:%d", m.From), string(b)).Err(); err != nil {
										log.Printf("redis publish ack error: %v", err)
									}
								}
							}
						}
					} else {
						// global broadcast
						for c := range h.clients {
							select {
							case c.send <- &m:
							default:
								close(c.send)
								delete(h.clients, c)
							}
						}
					}
				}
			}
		}
	}()
}

func (h *Hub) removeSub(channel string) {
	if global.GVA_REDIS == nil {
		return
	}
	h.subsMu.Lock()
	defer h.subsMu.Unlock()
	if s, ok := h.subs[channel]; ok {
		s.cancel()
		delete(h.subs, channel)
	}
}

func (h *Hub) publishToRedis(m *Message) {
	if global.GVA_REDIS == nil {
		return
	}
	var channel string
	if m.To != 0 {
		channel = fmt.Sprintf("user:%d", m.To)
	} else {
		channel = "broadcast"
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("redis marshal error: %v", err)
		return
	}
	if err := global.GVA_REDIS.Publish(context.Background(), channel, string(b)).Err(); err != nil {
		log.Printf("redis publish error: %v", err)
	}
}

func init() {
	go DefaultHub.Run()
}

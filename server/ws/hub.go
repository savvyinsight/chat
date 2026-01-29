package ws

import (
    "log"
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
}

func NewHub() *Hub {
    return &Hub{
        broadcast:  make(chan *Message, 256),
        register:   make(chan *Client, 128),
        unregister: make(chan *Client, 128),
        clients:    make(map[*Client]bool),
        users:      make(map[uint]map[*Client]bool),
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
                    }
                }
            }
            log.Printf("client unregistered: user=%d total=%d", c.userID, len(h.clients))
        case m := <-h.broadcast:
            if m.To != 0 {
                // targeted to a user
                if set, ok := h.users[m.To]; ok {
                    for c := range set {
                        select {
                        case c.send <- m:
                        default:
                            close(c.send)
                            delete(h.clients, c)
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

func init() {
    go DefaultHub.Run()
}

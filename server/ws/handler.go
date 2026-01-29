package ws

import (
	"net/http"
	"strconv"

	"chat/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS handles websocket requests from the peer.
// Expects a query parameter `user_id` (uint). JWT/auth integration will be added later.
func ServeWS(w http.ResponseWriter, r *http.Request) {
	// Try to authenticate via JWT token from Authorization header or "token" query
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	var userID uint
	if token != "" {
		user, err := service.AuthenticateToken(token)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		userID = user.ID
	} else {
		// fallback to user_id query for development
		userStr := r.URL.Query().Get("user_id")
		if userStr == "" {
			http.Error(w, "user_id or token required", http.StatusBadRequest)
			return
		}
		uid64, err := strconv.ParseUint(userStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		userID = uint(uid64)
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := NewClient(DefaultHub, conn, userID)
	DefaultHub.register <- client
	go client.WritePump()
	client.ReadPump()
}

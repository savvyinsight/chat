package ws

import (
	"testing"
	"time"
)

func TestHubDirectMessage(t *testing.T) {
	h := NewHub()
	go h.Run()

	c := NewClient(h, nil, 42)
	h.register <- c
	defer func() { h.unregister <- c }()

	// wait briefly for registration
	time.Sleep(10 * time.Millisecond)

	m := &Message{Type: "message", To: 42, Body: "hello"}
	h.broadcast <- m

	select {
	case got := <-c.send:
		if got.Body != "hello" {
			t.Fatalf("expected body 'hello', got '%s'", got.Body)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for direct message")
	}
}

func TestHubBroadcast(t *testing.T) {
	h := NewHub()
	go h.Run()

	c1 := NewClient(h, nil, 1)
	c2 := NewClient(h, nil, 2)
	h.register <- c1
	h.register <- c2
	defer func() { h.unregister <- c1; h.unregister <- c2 }()

	time.Sleep(10 * time.Millisecond)

	m := &Message{Type: "message", To: 0, Body: "hey all"}
	h.broadcast <- m

	for _, c := range []*Client{c1, c2} {
		select {
		case got := <-c.send:
			if got.Body != "hey all" {
				t.Fatalf("expected 'hey all', got '%s'", got.Body)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for broadcast message")
		}
	}
}

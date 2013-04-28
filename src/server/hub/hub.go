package hub

import "fmt"
import "server/document"

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Connections map[*Connection]bool

	// Inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	Register chan *Connection

	// Unregister requests from connections.
	Unregister chan *Connection

	Document: document.Document
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			fmt.Println("register",c)
			h.Connections[c] = true
		case c := <-h.Unregister:
			fmt.Println("unregister",c)
			delete(h.Connections, c)
			close(c.Send)
		case m := <-h.Broadcast:
			fmt.Println("broadcast",m)
			for c := range h.Connections {
				select {
				case c.Send <- m:
				default:
					close(c.Send)
					delete(h.Connections, c)
				}
			}
		}
	}
}

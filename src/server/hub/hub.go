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

}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type DocumentHub struct {
	// Registered connections.
	Connections map[*DocumentConnection]bool

	// Inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	Register chan *DocumentConnection

	// Unregister requests from connections.
	Unregister chan *DocumentConnection

	Document document.Document
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type ChatHub struct {
	// Registered connections.
	Connections map[*ChatConnection]bool

	// Inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	Register chan *ChatConnection

	// Unregister requests from connections.
	Unregister chan *ChatConnection

	Document document.Document
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			fmt.Println("register")
			h.Connections[c] = true
		case c := <-h.Unregister:
			fmt.Println("unregister")
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

//Make the call to update the document object and propigate the results to other users
func (h *DocumentHub) Run() {
	for {
		select {
		case c := <-h.Register:
			fmt.Println("Doc register")
			h.Connections[c] = true
		case c := <-h.Unregister:
			fmt.Println("Doc unregister")
			delete(h.Connections, c)
			close(c.Send)
		case m := <-h.Broadcast:
			fmt.Println("Doc broadcast",m)
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

//Pass the chats on to the right people
func (h *ChatHub) Run() {
	for {
		select {
		case c := <-h.Register:
			fmt.Println("Chat register")
			h.Connections[c] = true
		case c := <-h.Unregister:
			fmt.Println("Chat unregister")
			delete(h.Connections, c)
			close(c.Send)
		case m := <-h.Broadcast:
			fmt.Println("Chat broadcast",m)
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

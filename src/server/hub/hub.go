package hub

import "fmt"
import "server/document"
import "encoding/json"

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
//{"OpData":[{"Insert":"abc", "Position":0}],"Type":"Text","Name":"newDoc","Version":0,"Snapshot":"","Metadata":{"Creator":"","Ctime":"2013-04-29T14:59:36.346073-04:00","Mtime":"2013-04-29T14:59:36.346074-04:00","Sessions":{}}}
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
			//this is where ops should be processed and then sent
			//apparently this blocks? 
			var doc document.Document
        	err := json.Unmarshal(m, &doc)
        	if err != nil {
            	fmt.Println("Error:", err, string(m))
        	}
        	h.Document.ApplyOps(doc.OpData[0], doc.Version)
        	h.Document.BumpVersion()
        	json_bytes, _ := json.Marshal(h.Document)
			for c := range h.Connections {
				select {
				case c.Send <- json_bytes:
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

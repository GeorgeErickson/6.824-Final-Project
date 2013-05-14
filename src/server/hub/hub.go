package hub

import "fmt"
import "server/document"
import "encoding/json"
import "github.com/garyburd/redigo/redis"

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
	Broadcast chan Message

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

func (h *DocumentHub) Save(conn redis.Conn){
	json_bytes, _ := json.Marshal(h.Document)
	conn.Do("SET", h.Document.Name, string(json_bytes))
}

func getRedis() (redis.Conn, error) {
    c, err := redis.Dial("tcp", "127.0.0.1:6379")
    if err != nil {
        panic(err)
    }
    // if _, err := c.Do("AUTH", "davidgeorge"); err != nil {
    //      c.Close()
    //      return nil, err
    // }
    return c, err
}

//Make the call to update the document object and propigate the results to other users
//{"OpData":[[{"Insert":"abc", "Position":0}]],"Type":"Text","Name":"newDoc","Version":0,"Snapshot":"","Metadata":{"Creator":"","Ctime":"2013-04-29T14:59:36.346073-04:00","Mtime":"2013-04-29T14:59:36.346074-04:00","Sessions":{}}}
func (h *DocumentHub) Run() {
	redis_conn, _ := getRedis()
    //save document on close
    defer h.Save(redis_conn)
    defer redis_conn.Close()
	for {
		select {
		case c := <-h.Register:
			fmt.Println("Doc register")
			h.Connections[c] = true
			ndoc := document.Document{Name:h.Document.Name, Title:h.Document.Title, Version: h.Document.Version, Snapshot:h.Document.Snapshot}
			json_bytes, _ := json.Marshal(ndoc)
			c.Send <- Message{M:json_bytes}
		case c := <-h.Unregister:
			fmt.Println("Doc unregister")
			delete(h.Connections, c)
			close(c.Send)
		case mess := <-h.Broadcast:
			m := mess.M
			conn := mess.Conn
			//this is where ops should be processed and then sent
			//apparently this blocks? 
			var doc document.Document
        	err := json.Unmarshal(m, &doc)
        	if err != nil {
            	fmt.Println("Error:", err, string(m))
            	continue
        	}
        	fmt.Println(h.Document.Snapshot, doc)
        	error := h.Document.ApplyOps(doc.OpData[0], doc.Version)
        	if !error {
        		json_bytes, _ := json.Marshal(h.Document)
        		conn.Send <- Message{M:json_bytes}
        		continue
        	}
        	ndoc := document.Document{Name:h.Document.Name, ClientId:doc.ClientId, Title:h.Document.Title, Version: h.Document.Version, OpData:h.Document.OpData[len(h.Document.OpData)-(h.Document.Version-doc.Version):]}
        	json_bytes, _ := json.Marshal(ndoc)
		    if(h.Document.Version % 20 == 0){
		    	h.Save(redis_conn)
		    }
			for c := range h.Connections {
				if c == conn{
					continue
				}
				select {
				case c.Send <- Message{M:json_bytes}:
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

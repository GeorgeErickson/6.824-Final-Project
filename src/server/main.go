package main

import (
    "github.com/hoisie/web"
    "github.com/garyburd/go-websocket/websocket"
    "github.com/garyburd/redigo/redis"
    "html/template"
    "bytes"
    "log"
    "strconv"
    "fmt"
    "encoding/json"
    "io"
    "server/hub"
)

/*
ToDO:
Each document needs its own hub to manage connections
need a global mapping of documents to hubs, each one needs to run
modify hub to handle documents and edits
copy for chats

should make seperate document hubs and chat hubs
*/
func parseTemplate(file string, data interface{}) (out []byte, error error) {
        var buf bytes.Buffer
        t, err := template.ParseFiles(file)
        if err != nil {
                return nil, err
        }
        err = t.Execute(&buf, data)
        if err != nil {
                return nil, err
        }
        return buf.Bytes(), nil
}

// serverWs handles websocket requests from the client.
func serveWs(ctx *web.Context) {
	w := ctx.ResponseWriter
	r := ctx.Request
	if r.Header.Get("Origin") != "http://"+r.Host {
		ctx.Abort(403, "Origin not allowed")
		// http.Error(w, "Origin not allowed", 403)
		return
	}
	ws, err := websocket.Upgrade(w, r.Header, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		ctx.Abort(400, "Not a websocket handshake")
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	c := &hub.Connection{Send: make(chan []byte, 256), Ws: ws, H:h}
	h.Register <- c
	go c.WritePump()
	c.ReadPump()
}

// documentStream handles websocket requests from the client for a particular document
func documentStream(ctx *web.Context, documentId string) {
    w := ctx.ResponseWriter
    r := ctx.Request
    if r.Header.Get("Origin") != "http://"+r.Host {
        ctx.Abort(403, "Origin not allowed")
        // http.Error(w, "Origin not allowed", 403)
        return
    }
    ws, err := websocket.Upgrade(w, r.Header, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); ok {
        ctx.Abort(400, "Not a websocket handshake")
        return
    } else if err != nil {
        log.Println(err)
        return
    }
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(err)
    }
    defer c.Close()
    s, err := redis.String(c.Do("GET", documentId))
    //make new document at that stringId
    if err {
        doc := document.newDoc(documentId)
        c := &hub.Connection{Send: make(chan []byte, 256), Ws: ws, H:h}
        document_hubs[doc.Name] = hub.Hub{
                                    Document:    doc,
                                    Broadcast:   make(chan []byte),
                                    Register:    make(chan *hub.Connection),
                                    Unregister:  make(chan *hub.Connection),
                                    Connections: make(map[*hub.Connection]bool),
                                }
        document_hubs[doc.Name].Register <- c
        go c.WritePump()
        c.ReadPump()
    } else{
        var doc document.Document
        err := json.Unmarshal([]byte(s), &doc)
        if err != nil {
            fmt.Println("error:", err)
        }
        c := &hub.Connection{Send: make(chan []byte, 256), Ws: ws, H:h}
        document_hubs[doc.Name].Register <- c
        go c.WritePump()
        c.ReadPump()
    }
    
}

// serverWs handles websocket requests from the client for a particular chat
func chatStream(ctx *web.Context, chatId string) {
    w := ctx.ResponseWriter
    r := ctx.Request
    if r.Header.Get("Origin") != "http://"+r.Host {
        ctx.Abort(403, "Origin not allowed")
        // http.Error(w, "Origin not allowed", 403)
        return
    }
    ws, err := websocket.Upgrade(w, r.Header, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); ok {
        ctx.Abort(400, "Not a websocket handshake")
        return
    } else if err != nil {
        log.Println(err)
        return
    }
    c := &hub.Connection{Send: make(chan []byte, 256), Ws: ws, H:h}
    h.Register <- c
    go c.WritePump()
    c.ReadPump()
}


func home(ctx *web.Context) { 
    home, error := parseTemplate("home.html", map[string]string{"Url": ctx.Server.Config.Addr, "Port": strconv.Itoa(ctx.Server.Config.Port)})
    var buf bytes.Buffer
    if error != nil {
    	buf.Write([]byte(error.Error()))
    	io.Copy(ctx, &buf)
    }
    buf.Write(home)
    io.Copy(ctx, &buf)
} 

func getDocuments(ctx *web.Context){
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(err)
    }
    defer c.Close()
    s, err := redis.Strings(c.Do("KEYS", "*"))
    fmt.Printf("%#v\n", s)
    var buf bytes.Buffer
    for _,str := range s {
        buf.Write([]byte(str))
    }
    io.Copy(ctx, &buf)
}

var mainHub = hub.Hub{
		Broadcast:   make(chan []byte),
		Register:    make(chan *hub.Connection),
		Unregister:  make(chan *hub.Connection),
		Connections: make(map[*hub.Connection]bool),
	}

func setup() map[string]Hub {
    document_hubs := map[string]Hub{}
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(err)
    }
    defer c.Close()
    s, err := redis.Strings(c.Do("KEYS", "*"))
    var documents map[string]Document
    for _,str := range s {
        jdoc, err := redis.String(c.Do("GET", str))
        var doc Document
        err := json.Unmarshal([]byte(jdoc), &doc)
        if err != nil {
            fmt.Println("error:", err)
        }
        documents[doc.Name] = doc
    }
    for _, doc := range documents {
        var h = hub.Hub{
            Document:    doc,
            Broadcast:   make(chan []byte),
            Register:    make(chan *hub.Connection),
            Unregister:  make(chan *hub.Connection),
            Connections: make(map[*hub.Connection]bool),
        }
        go h.Run()
        document_hubs[doc.name] = h
    }
    return document_hubs
}

func main() {
    go mainHub.Run()
    document_hubs = setup()
	web.Config.Addr = "127.0.0.1"
	web.Config.Port = 8000
    web.Get("/", home)
    web.Get("/ws", serveWs)
    web.Get("/rest/documents", getDocuments)
    web.Get("/documents/(.*)/", documentStream)
    web.Get("/chat/(.*)/", chatStream)
    web.Run("127.0.0.1:8000")
}
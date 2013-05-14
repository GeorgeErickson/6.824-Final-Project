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
    "io/ioutil"
    "server/hub"
    "server/document"
    "os"
    "strings"
)

/*
ToDO:
Integrate messaging protocol and document change protocol
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
	c := &hub.Connection{Send: make(chan []byte, 256), Ws: ws, H:mainHub}
	mainHub.Register <- c
	go c.WritePump()
	c.ReadPump()
}

func setDocumentTitle(ctx *web.Context, documentId string){
    redis_conn, _ := getRedis()
    defer redis_conn.Close()
    
    //make new document at that stringId
    if _, ok := document_hubs[documentId]; ok {
        body, _ := ioutil.ReadAll(ctx.Request.Body)
        var ndoc document.Document
        json.Unmarshal(body, &ndoc)
        h := document_hubs[documentId]
        h.Document.Title = ndoc.Title
        var buf bytes.Buffer
        doc := document.Document{Name:h.Document.Name, Title:h.Document.Title, Version: h.Document.Version}
        h.Save(redis_conn)
        doc.ClientId = ndoc.ClientId
        json_bytes, _ := json.Marshal(doc)
        mainHub.Broadcast <- json_bytes
        buf.Write(json_bytes)
        io.Copy(ctx, &buf)
    } else{
        ctx.Abort(404, "Document does not exist")
        return
    }
}

func deleteDocument(ctx *web.Context, documentId string){
    redis_conn, _ := getRedis()
    defer redis_conn.Close()
    
    //make new document at that stringId
    if _, ok := document_hubs[documentId]; ok {
        h := document_hubs[documentId]
        var buf bytes.Buffer
        doc := document.Document{Name:h.Document.Name, Title:h.Document.Title, Version: h.Document.Version}
        redis_conn.Do("DEL", documentId)
        json_bytes, _ := json.Marshal(doc)
        h.Broadcast <- hub.Message{M:json_bytes}
        buf.Write(json_bytes)
        io.Copy(ctx, &buf)
    } else{
        ctx.Abort(404, "Document does not exist")
        return
    }
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
    redis_conn, _ := getRedis()
    defer redis_conn.Close()
    s, err := redis.String(redis_conn.Do("GET", documentId))
    //make new document at that stringId
    if err != nil {
        doc := document.NewDoc(documentId)
        document_hubs[doc.Name] = hub.DocumentHub{
                                    Document:    doc,
                                    Broadcast:   make(chan hub.Message),
                                    Register:    make(chan *hub.DocumentConnection),
                                    Unregister:  make(chan *hub.DocumentConnection),
                                    Connections: make(map[*hub.DocumentConnection]bool),
                                }
        h := document_hubs[doc.Name]
        go h.Run()
        json_bytes, _ := json.Marshal(doc)
        redis_conn.Do("SET", documentId, string(json_bytes))
        c := &hub.DocumentConnection{Send: make(chan hub.Message, 256), Ws: ws, H:document_hubs[doc.Name]}
        document_hubs[doc.Name].Register <- c
        mainHub.Broadcast <- json_bytes
        go c.WritePump()
        c.ReadPump()
    } else{
        var doc document.Document
        fmt.Println(s)
        err := json.Unmarshal([]byte(s), &doc)
        if err != nil {
            fmt.Println("Error:", err, s)
        }
        c := &hub.DocumentConnection{Send: make(chan hub.Message, 256), Ws: ws, H:document_hubs[doc.Name]}
        document_hubs[doc.Name].Register <- c
        go c.WritePump()
        c.ReadPump()
    }
}

// serverWs handles websocket requests from the client for a particular chat
func chatStream(ctx *web.Context, documentId string) {
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
    redis_conn, _ := getRedis()
    defer redis_conn.Close()
    s, err := redis.String(redis_conn.Do("GET", documentId))
    //make new document at that stringId
    if err != nil {
        doc := document.NewDoc(documentId)
        chat_hubs[doc.Name] = hub.ChatHub{
                                    Document:    doc,
                                    Broadcast:   make(chan []byte),
                                    Register:    make(chan *hub.ChatConnection),
                                    Unregister:  make(chan *hub.ChatConnection),
                                    Connections: make(map[*hub.ChatConnection]bool),
                                }
        h := chat_hubs[doc.Name]
        go h.Run()
        c := &hub.ChatConnection{Send: make(chan []byte, 256), Ws: ws, H:chat_hubs[doc.Name]}
        chat_hubs[doc.Name].Register <- c
        go c.WritePump()
        c.ReadPump()
    } else{
        var doc document.Document
        err := json.Unmarshal([]byte(s), &doc)
        if err != nil {
            fmt.Println("error:", err)
        }
        c := &hub.ChatConnection{Send: make(chan []byte, 256), Ws: ws, H:chat_hubs[doc.Name]}
        chat_hubs[doc.Name].Register <- c
        go c.WritePump()
        c.ReadPump()
    }
}


func home(ctx *web.Context) { 
    home, error := parseTemplate("client/public/index.html", map[string]string{"Url": ctx.Server.Config.Addr, "Port": strconv.Itoa(ctx.Server.Config.Port)})
    var buf bytes.Buffer
    if error != nil {
    	buf.Write([]byte(error.Error()))
    	io.Copy(ctx, &buf)
    }
    buf.Write(home)
    io.Copy(ctx, &buf)
} 

func getDocuments(ctx *web.Context){
    c, _ := getRedis()
    defer c.Close()
    s, _ := redis.Strings(c.Do("KEYS", "*"))
    fmt.Printf("%#v\n", s)
    var documents = map[string]document.Document{}
    for _,str := range s {
        jdoc, err := redis.String(c.Do("GET", str))
        var doc document.Document
        error := json.Unmarshal([]byte(jdoc), &doc)
        if error != nil {
            fmt.Println("Document setup error:", err)
        }
        documents[doc.Name] = doc
        fmt.Println("Document: ", doc.Name)
    }

    var buf bytes.Buffer
    json_bytes, _ := json.Marshal(documents)
    buf.Write(json_bytes)
    io.Copy(ctx, &buf)
}

var mainHub = hub.Hub{
		Broadcast:   make(chan []byte),
		Register:    make(chan *hub.Connection),
		Unregister:  make(chan *hub.Connection),
		Connections: make(map[*hub.Connection]bool),
	}

func getRedis() (redis.Conn, error) {
    c, err := redis.Dial("tcp", "pub-redis-11830.us-east-1-4.1.ec2.garantiadata.com:11830")
    if err != nil {
        panic(err)
    }
    if _, err := c.Do("AUTH", "davidgeorge"); err != nil {
         c.Close()
         return nil, err
    }
    return c, err
}

func setupDocuments() map[string]hub.DocumentHub {
    d_hubs := map[string]hub.DocumentHub{}
    c, _ := getRedis()
    defer c.Close()
    s, _ := redis.Strings(c.Do("KEYS", "*"))
    var documents = map[string]document.Document{}
    for _,str := range s {
        jdoc, err := redis.String(c.Do("GET", str))
        var doc document.Document
        error := json.Unmarshal([]byte(jdoc), &doc)
        if error != nil {
            fmt.Println("Document setup error:", err)
        }
        documents[doc.Name] = doc
        fmt.Println("Document: ", doc.Name)
    }
    for _, doc := range documents {
        var h = hub.DocumentHub{
            Document:    doc,
            Broadcast:   make(chan hub.Message),
            Register:    make(chan *hub.DocumentConnection),
            Unregister:  make(chan *hub.DocumentConnection),
            Connections: make(map[*hub.DocumentConnection]bool),
        }
        go h.Run()
        d_hubs[doc.Name] = h
    }
    return d_hubs
}

func setupChats() map[string]hub.ChatHub {
    d_hubs := map[string]hub.ChatHub{}
    c, _ := getRedis()
    defer c.Close()
    s, _ := redis.Strings(c.Do("KEYS", "*"))
    var documents = map[string]document.Document{}
    for _,str := range s {
        jdoc, err := redis.String(c.Do("GET", str))
        var doc document.Document
        error := json.Unmarshal([]byte(jdoc), &doc)
        if error != nil {
            fmt.Println("chat setup error:", err)
        }
        documents[doc.Name] = doc
    }
    for _, doc := range documents {
        var h = hub.ChatHub{
            Document:    doc,
            Broadcast:   make(chan []byte),
            Register:    make(chan *hub.ChatConnection),
            Unregister:  make(chan *hub.ChatConnection),
            Connections: make(map[*hub.ChatConnection]bool),
        }
        go h.Run()
        d_hubs[doc.Name] = h
    }
    return d_hubs
}

var document_hubs = map[string]hub.DocumentHub{}
var chat_hubs = map[string]hub.ChatHub{}

func main() {
    go mainHub.Run()
    document_hubs = setupDocuments()
    chat_hubs = setupChats()
    // Get hostname if specified
    args := os.Args;
    host := "0.0.0.0:8080"
    if len(args) == 2 {
        host = args[1]
    }
    parts := strings.Split(host, ":")
    fmt.Println(parts)
	web.Config.Addr = parts[0]
	web.Config.Port, _ = strconv.Atoi(parts[1])
    web.Config.StaticDir = "client/public/"
    web.Get("/", home)
    web.Get("/ws", serveWs)
    web.Get("/rest/documents", getDocuments)
    web.Put("/rest/documents/(.*)", setDocumentTitle)
    web.Delete("/rest/documents/(.*)", deleteDocument)
    web.Get("/documents/(.*)", documentStream)
    web.Get("/chat/(.*)", chatStream)
    web.Run(host)

}
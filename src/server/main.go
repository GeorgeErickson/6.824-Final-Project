package main

import (
    "github.com/hoisie/web"
    "github.com/garyburd/go-websocket/websocket"
    "html/template"
    "bytes"
    "log"
    "strconv"
    "encoding/json"
    "io"
    "server/hub"
)

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

//return the html for a document
func documentSetup(ctx *web.Context, documentName string) { 
    home, error := parseTemplate("home.html", map[string]string{"Url": ctx.Server.Config.Addr, "Port": strconv.Itoa(ctx.Server.Config.Port)})
    var buf bytes.Buffer
    if error != nil {
        buf.Write([]byte(error.Error()))
        io.Copy(ctx, &buf)
    }
    buf.Write(home)
    io.Copy(ctx, &buf)
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

var h = hub.Hub{
		Broadcast:   make(chan []byte),
		Register:    make(chan *hub.Connection),
		Unregister:  make(chan *hub.Connection),
		Connections: make(map[*hub.Connection]bool),
	}

func main() {
	go h.Run()
	web.Config.Addr = "127.0.0.1"
	web.Config.Port = 8000
    web.Get("/", home)
    web.Get("/ws", serveWs)
    web.Get("/(.*)/", documentSetup)
    web.Run("127.0.0.1:8000")
}
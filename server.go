package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Conn representa uma conexão do websocket
type Conn struct {
	WS   *websocket.Conn
	Send chan string
}

// SendToHub ...
func (conn *Conn) SendToHub() {
	defer conn.WS.Close()
	for {
		_, msg, err := conn.WS.ReadMessage()
		if err != nil {
			return
		}
		DefaultHub.Echo <- string(msg)
	}
}

// Write ...
func (conn *Conn) Write(msg string) error {
	return conn.WS.WriteMessage(websocket.TextMessage, []byte(msg))
}

// ReceiveFromHub ...
func (conn *Conn) ReceiveFromHub() {
	defer conn.WS.Close()
	for {
		conn.Write(<-conn.Send)
	}
}

// WSHandler ...
func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn := &Conn{
		Send: make(chan string),
		WS:   ws,
	}

	DefaultHub.Join <- conn

	go conn.SendToHub()
	conn.ReceiveFromHub()
}

// WSUpgrader ...
var WSUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

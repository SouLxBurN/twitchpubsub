package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// NewConnection
func NewConnection(connectionUrl string) *websocket.Conn {
	// Connect to Websocket
	c, _, err := websocket.DefaultDialer.Dial(connectionUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

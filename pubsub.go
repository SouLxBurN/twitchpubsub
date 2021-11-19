package main

import (
	"log"
	"os"
	"soulxbot/process"
	"soulxbot/schema"

	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	dotenv "github.com/joho/godotenv"
)

const defaultWebsocketURL = "wss://pubsub-edge.twitch.tv"

func main() {
	if err := dotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	authToken := os.Getenv("SOULXBOT_AUTHTOKEN")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Connect to Websocket
	c, _, err := websocket.DefaultDialer.Dial(defaultWebsocketURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Subscribe to Topics
	listen := schema.ListenRequest{
		Type: "LISTEN",
		Data: schema.ListenRequestData{
			Topics:    []string{"channel-points-channel-v1.31568083"},
			AuthToken: authToken,
		},
	}
	if err := c.WriteJSON(listen); err != nil {
		log.Fatal("Failed to subscribe to Topics", err)
	}

	// ReadMessage Process
	done := make(chan struct{})
	go process.ReadMessage(c, done)
	Listen(c, interrupt, done)
}

func Listen(c *websocket.Conn, interrupt chan os.Signal, done chan struct{}) {
	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := c.WriteJSON(schema.ListenRequest{Type: "PING"})
			if err != nil {
				log.Println("Ping Failed:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

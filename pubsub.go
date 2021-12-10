package main

import (
	"github.com/soulxburn/twitchpubsub/process"
	"github.com/soulxburn/twitchpubsub/schema"
	"log"
	"os"

	"os/signal"

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
	process.Listen(c, interrupt, done)
}

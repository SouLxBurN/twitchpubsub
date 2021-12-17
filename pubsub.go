package main

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/soulxburn/twitchpubsub/process"
	"github.com/soulxburn/twitchpubsub/schema"
	"github.com/soulxburn/twitchpubsub/twitch"
	"github.com/soulxburn/twitchpubsub/ws"

	"os/signal"

	dotenv "github.com/joho/godotenv"
)

const TWITCH_WS_URL = "wss://pubsub-edge.twitch.tv"

func main() {
	if err := dotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	clientID := os.Getenv("SOULXBOT_CLIENTID")
	clientSecret := os.Getenv("SOULXBOT_CLIENTSECRET")
	authToken := os.Getenv("SOULXBOT_AUTHTOKEN")
	refreshToken := os.Getenv("SOULXBOT_REFRESHTOKEN")
	auth := twitch.NewAuthTokenProxy(clientID, clientSecret, authToken, refreshToken)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		aToken, err := auth.GetAuthToken()
		if err != nil {
			log.Fatal("Failed to retrieve/refresh auth token: ", err)
		}

		c := ws.NewConnection(TWITCH_WS_URL)
		// Subscribe to Topics
		listen := schema.ListenRequest{
			Type: "LISTEN",
			Data: schema.ListenRequestData{
				Topics:    []string{"channel-points-channel-v1.31568083"},
				AuthToken: aToken,
			},
		}
		if err := c.WriteJSON(listen); err != nil {
			log.Fatal("Failed to subscribe to Topics", err)
		}

		// ReadMessage Process
		done := make(chan struct{})
		go process.ReadMessage(c, done)
		go process.StartPinging(c, done)

		select {
		case <-done:
		case <-interrupt:
			defer c.Close()
			log.Println("interrupt received")

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

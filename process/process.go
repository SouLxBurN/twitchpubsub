package process

import (
	"encoding/json"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/soulxburn/twitchpubsub/schema"
)

const (
	SELECT_BACKGROUND_RWD_ID = "6f111bb7-315a-462a-a4e0-54333ca06a14"
	RANDOM_BACKGROUND_RWD_ID = "2b43e9da-b7db-4932-8a05-67422a93d1f8"
)

var currentBackground int = 0

// ReadMessage
func ReadMessage(c *websocket.Conn, done chan struct{}) {
	defer close(done)
	rand.Seed(time.Now().UnixNano())
	for {
		mtype, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error:", err)
			return
		}
		var messageObject = &schema.Message{}
		if err := json.Unmarshal(message, messageObject); err != nil {
			log.Println("Failed to Unmarshal: ", err)
		}
		switch messageObject.Type {
		case "RESPONSE":
			if messageObject.Error == "ERR_BADAUTH" {
				log.Println("Authorization Failed.")
				time.Sleep(5 * time.Second)
				return
			}
		case "RECONNECT":
			log.Println("Received Reconnect Notice")
			return
		case "MESSAGE":
			var tpm = &schema.TwitchPubMessage{}
			if err := json.Unmarshal([]byte(messageObject.Data.Message), tpm); err != nil {
				log.Println("Failed to Unmarshal Twitch Message: ", err)
			}

			if tpm.Type == "reward-redeemed" {
				log.Println("Reward ID: ", tpm.Data.Redemption.ID)
				switch tpm.Data.Redemption.Reward.ID {
				case SELECT_BACKGROUND_RWD_ID:
					input, err := strconv.Atoi(tpm.Data.Redemption.UserInput[:2])
					if err != nil {
						if err := changeBackground(input); err != nil {
							log.Println("Error while selecting background: ", err)
						}
					}
				case RANDOM_BACKGROUND_RWD_ID:
					input := currentBackground
					for ; input == currentBackground; input = (rand.Int() % (len(wallpapers) - 1)) + 1 {
					}
					if err := changeBackground(input); err != nil {
						log.Println("Error while randomizing background: ", err)
					}
				}
			}
		}
		log.Printf("recv(%v): %s", mtype, message)
	}
}

// StartPinging
func StartPinging(c *websocket.Conn, done chan struct{}) {
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
		case <-done:
			return
		}
	}
}

// changeBackground
func changeBackground(wpIdx int) error {
	wallpaper := wallpapers[strconv.Itoa(wpIdx)]
	if wallpaper != "" {
		cmd := exec.Command("feh", "--bg-scale", wallpaper)
		stdout, err := cmd.Output()
		log.Println("feh Output: ", stdout)
		currentBackground = wpIdx
		return err
	}
	return nil
}

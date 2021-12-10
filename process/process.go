package process

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
	"github.com/soulxburn/twitchpubsub/schema"
)

const BACKGROUND_REWARD_ID = "6f111bb7-315a-462a-a4e0-54333ca06a14"

var wallpapers = map[string]string{
	"1": "/home/soulxburn/.config/i3/digital-landscape-wallpaper.jpg",
	"2": "/home/soulxburn/wallpaper/The_Pink_Sunset_Wallpaper_2560x1600.jpg",
	"3": "/home/soulxburn/wallpaper/black-and-red-abstract-painting.jpeg",
	"4": "/home/soulxburn/wallpaper/3399647.jpeg",
	"5": "/home/soulxburn/wallpaper/95-952362_dark-wood-wallpapers-terminal-background-image-dark.jpg",
}

// ReadMessage
func ReadMessage(c *websocket.Conn, done chan struct{}) {
	defer close(done)
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

		if messageObject.Type == "MESSAGE" {
			var tpm = &schema.TwitchPubMessage{}
			if err := json.Unmarshal([]byte(messageObject.Data.Message), tpm); err != nil {
				log.Println("Failed to Unmarshal Twitch Message: ", err)
			}

			if tpm.Type == "reward-redeemed" {
				log.Println("Reward ID: ", tpm.Data.Redemption.ID)
				if tpm.Data.Redemption.Reward.ID == BACKGROUND_REWARD_ID {
					wallpaper := wallpapers[tpm.Data.Redemption.UserInput[:1]]
					if wallpaper != "" {
						cmd := exec.Command("feh", "--bg-scale", wallpaper)
						stdout, err := cmd.Output()
						log.Println("feh Output: ", stdout, err)
					}
				}
			}
		}
		log.Printf("recv(%v): %s", mtype, message)
	}
}

// Listen
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

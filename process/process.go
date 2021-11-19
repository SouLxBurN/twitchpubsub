package process

import (
	"encoding/json"
	"log"
	"os/exec"
	"soulxbot/schema"

	"github.com/gorilla/websocket"
)

const BACKGROUND_REWARD_ID = "6f111bb7-315a-462a-a4e0-54333ca06a14"

var wallpapers = map[string]string{
	"1": "/home/soulxburn/.config/i3/digital-landscape-wallpaper.jpg",
	"2": "/home/soulxburn/wallpaper/The_Pink_Sunset_Wallpaper_2560x1600.jpg",
	"3": "/home/soulxburn/wallpaper/black-and-red-abstract-painting.jpeg",
	"4": "/home/soulxburn/wallpaper/3399647.jpeg",
	"5": "/home/soulxburn/wallpaper/95-952362_dark-wood-wallpapers-terminal-background-image-dark.jpg",
}

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

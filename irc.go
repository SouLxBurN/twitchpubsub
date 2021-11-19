package main

import (
	"fmt"
	"log"
	"os"

	"os/exec"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v2"
	dotenv "github.com/joho/godotenv"
)

func main() {
	if err := dotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	user := os.Getenv("SOULXBOT_USER")
	oauth := os.Getenv("SOULXBOT_OAUTH")

	firstMessage := false
	client := twitch.NewClient(user, oauth)

	client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		fmt.Printf("Notice: %s\n", message.Message)
	})

	wallpapers := map[string]string{
		"1": "/home/soulxburn/.config/i3/digital-landscape-wallpaper.jpg",
		"2": "/home/soulxburn/wallpaper/The_Pink_Sunset_Wallpaper_2560x1600.jpg",
		"3": "/home/soulxburn/wallpaper/black-and-red-abstract-painting.jpeg",
		"4": "/home/soulxburn/wallpaper/3399647.jpeg",
		"5": "/home/soulxburn/wallpaper/95-952362_dark-wood-wallpapers-terminal-background-image-dark.jpg",
	}

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if firstMessage {
			firstMessage = false
			client.Say(message.Channel, fmt.Sprintf("Congratulations %s! You're first!", message.User.DisplayName))
		}

		fmt.Printf("%s: %s\n", message.User.DisplayName, message.Message)

		if strings.Contains(message.Tags["msg-id"], "highlighted-message") {
			wallpaper := wallpapers[message.Message[:1]]
			if wallpaper != "" {
				cmd := exec.Command("feh", "--bg-scale", wallpaper)
				stdout, err := cmd.Output()
				fmt.Println(stdout, err)
			}
		}
	})

	client.Join("SouLxBurN")

	if err := client.Connect(); err != nil {
		panic(err)
	}

}

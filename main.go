package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"team-radio/auth"
	"team-radio/utils"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func main() {
	godotenv.Load(".env")
	var prevSongLink string

	botToken := os.Getenv("SLACK_BOT_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_AUTH_TOKEN")

	code, err := auth.FetchAuthCode()

	if err != nil {
		log.Println(err)
		return
	}

	accessToken, _, err := auth.GetTokens(code)

	log.Printf("Token: %s", accessToken)

	if err != nil {
		log.Println(err)
		return
	}

	bot := slacker.NewClient(botToken, appToken)

	go utils.LogCommandEvents(bot.CommandEvents())

	bot.Command("spin", &slacker.CommandDefinition{
		Description: "Add music to team radio",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			if prevSongLink == botCtx.Event().Text {
				return
			}

			prevSongLink = botCtx.Event().Text

			var tracksToAdd []string

			playlistId := os.Getenv("PLAYLIST_ID")

			songId, err := utils.PrepareSongID(botCtx.Event().Text)
			if err != nil {
				response.Reply("Sorry, the link is not from spotify.")
				log.Printf("Error getting song ID: %v", err)
				return
			}

			log.Printf("new music link: %v | song ID: %v | playlist ID: %v ", botCtx.Event().Text, songId, playlistId)

			tracksToAdd = append(tracksToAdd, songId)

			code := utils.AddMusicToPlaylist(accessToken, tracksToAdd, playlistId)

			if code != 201 {
				log.Println("could not add music to playlist")
				response.Reply("Sorry, couldnt add song to playlist")
				return
			}

			response.Reply(fmt.Sprintf("%s put something on the turntable", botCtx.Event().UserProfile.FirstName))
		},
	})

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	err = bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

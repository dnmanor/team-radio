package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/shomali11/slacker"
)

func AddMusicToPlaylist(token string, uris []string, playlistID string) int {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
	method := "POST"

	payload := map[string]interface{}{
		"uris":     uris,
		"position": 0,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return 500
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return 500
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return resp.StatusCode
	}
	defer resp.Body.Close()

	bodyOfRes, _ := ioutil.ReadAll(resp.Body)

	// Print the response status code & Body(snapshot_id if OK)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(bodyOfRes))

	return resp.StatusCode
}

func LogCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		log.Println("Command Events Log")
		log.Printf("cmd: %s| @ %s | user: %s", event.Command, event.Timestamp, event.Event.UserID)
	}

}

func PrepareSongID(musicLink string) (string, error) {
	isSpotifyLink := strings.Contains(musicLink, "spotify")

	if !isSpotifyLink {
		return "", fmt.Errorf("link is not from spotify")
	}

	linkSplit := strings.Split(musicLink, "/")
	rawSongIdWithTracking := strings.TrimRight(linkSplit[len(linkSplit)-1], ">")
	rawSongId := strings.Split(rawSongIdWithTracking, "?")[0]

	songId := fmt.Sprintf("spotify:track:%s", rawSongId)

	return songId, nil
}

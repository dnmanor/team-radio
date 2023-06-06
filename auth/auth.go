package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
)

func extractCodeValue(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	values, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		return "", err
	}

	codeValue := values.Get("code")
	return codeValue, nil
}

//todo: handle auth state better using the spotify state feature rather than a fixed state
func FetchAuthCode() (string, error) {
	godotenv.Load("../.env")

	authURL := os.Getenv("SPOTIFY_AUTHCODE_URL")
	username := os.Getenv("SPOTIFY_USERNAME")
	password := os.Getenv("SPOTIFY_PASSWORD")
	WAIT_TIME := 7000

	page := rod.New().MustConnect().MustPage(authURL)

	page.MustElement("#login-username").MustInput(username)
	page.MustElement("#login-password").MustInput(password)
	page.MustElement("#login-button").MustClick()
	page.MustWaitNavigation()
	// page.MustWaitLoad()
	time.Sleep(time.Millisecond * time.Duration(WAIT_TIME)) //interchangeable with the line above if you got good internet

	//comment out any time you change the access scopes and need to re-authenticate
	// urlContainsAuthCode := strings.Contains(page.MustInfo().URL, "callback?code=")

	// if !urlContainsAuthCode {
	// 	page.MustElement("[data-testid=\"auth-accept\"]").MustClick()
	// 	page.MustWaitLoad()
	// 	time.Sleep(time.Millisecond * time.Duration(WAIT_TIME))
	// 	page.MustScreenshotFullPage("ce.png")
	// }

	URL := page.MustInfo().URL

	code, err := extractCodeValue(URL)
	if err != nil {
		log.Println("Error:", err)
	}

	page.MustClose()

	return code, nil
}

func GetTokens(code string) (string, string, error) {
	godotenv.Load("../.env")

	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_SECRET")
	scopes := "user-read-email user-read-private playlist-modify-private playlist-modify-public playlist-read-private"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", spotifyClientID)
	data.Set("client_secret", spotifyClientSecret)
	data.Set("scope", scopes)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to get access token: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}

	var tokenResp map[string]interface{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	accessToken := tokenResp["access_token"].(string)
	refreshToken := tokenResp["refresh_token"].(string)

	return accessToken, refreshToken, nil

}

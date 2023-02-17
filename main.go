package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type TwitchUser struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}

type TwitchStream struct {
	ViewerCount int `json:"viewer_count"`
}

type TwitchTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	CreatedAt    int64  `json:"created_at"`
}

func main() {
	// Load Twitch app credentials from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Twitch access token using client credentials grant flow
	access_token, err := getTwitchAccessToken()
	if err != nil {
		log.Fatal(err)
	}

	// Get Twitch user information using access token
	user, err := getTwitchUser(access_token, "vladtenten") // Replace with your Twitch username
	if err != nil {
		log.Fatal(err)
	}

	// Get Twitch stream information using access token and user ID
	stream, err := getTwitchStream(access_token, user.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Print viewer count
	if stream == nil {
		log.Println("User is currently offline")
	} else {
		log.Println("Viewer Count:", stream.ViewerCount)
	}
}

func getTwitchAccessToken() (string, error) {
	data := url.Values{}

	data.Set("grant_type", "client_credentials")
	data.Set("client_id", os.Getenv("TWITCH_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("TWITCH_CLIENT_SECRET"))

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp TwitchTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func getTwitchUser(access_token string, username string) (*TwitchUser, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
	req.Header.Set("Authorization", "Bearer "+access_token)

	q := req.URL.Query()
	q.Add("login", username)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userResp struct {
		Data []TwitchUser `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		return nil, err
	}

	if len(userResp.Data) == 0 {
		return nil, nil
	}

	user := userResp.Data[0]

	return &user, nil
}

func getTwitchStream(access_token string, user_id string) (*TwitchStream, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
	req.Header.Set("Authorization", "Bearer "+access_token)

	q := req.URL.Query()
	q.Add("user_id", user_id)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var streamResp struct {
		Data []TwitchStream `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&streamResp)
	if err != nil {
		return nil, err
	}

	if len(streamResp.Data) == 0 {
		return nil, nil
	}

	stream := streamResp.Data[0]

	return &stream, nil
}

func refreshTwitchAccessToken(tokenResp TwitchTokenResponse) (string, error) {
	now := time.Now().Unix()
	if now >= tokenResp.CreatedAt+int64(tokenResp.ExpiresIn) {
		data := url.Values{}

		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", tokenResp.RefreshToken)
		data.Set("client_id", os.Getenv("TWITCH_CLIENT_ID"))
		data.Set("client_secret", os.Getenv("TWITCH_CLIENT_SECRET"))

		req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var newTokenResp TwitchTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&newTokenResp)
		if err != nil {
			return "", err
		}

		return newTokenResp.AccessToken, nil
	}
	return tokenResp.AccessToken, nil
}

package auth0

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"fmt"
	"github.com/joho/godotenv"
	
)

func GetAccessToken() (string, error) {
	godotenv.Load()
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	audience := os.Getenv("AUTH0_AUDIENCE")
	url := os.Getenv("AUTH0_TOKEN_URL")
	if clientID == "" || clientSecret == "" || audience == "" || url == "" {
		return "", fmt.Errorf("missing required environment variables")
	}

	payload := strings.NewReader(`{
		"client_id":"` + clientID + `",
		"client_secret":"` + clientSecret + `",
		"audience":"` + audience + `",
		"grant_type":"client_credentials"
	}`)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}
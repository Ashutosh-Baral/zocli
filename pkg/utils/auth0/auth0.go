package auth0

import (
    "encoding/json"
    "io"
    "net/http"
    "strings"
)

func GetAccessToken() (string, error) {
    url := "https://dev-qjq8cszju2kp1lxu.us.auth0.com/oauth/token"
    payload := strings.NewReader(`{
        "client_id":"1vWvsjV50UDVzKv7KCfuqp7Ubtwn2qew",
        "client_secret":"jtBtsgjtP2LaA0fwYjMcpfxFVmZu7NjOCsQ-jfvH3Owotdpl0ZYXs6GFo2Q7edm9",
        "audience":"https://dev-qjq8cszju2kp1lxu.us.auth0.com/api/v2/",
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
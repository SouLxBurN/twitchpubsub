package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	TWITCH_OAUTH_API = "https://id.twitch.tv/oauth2"
	TOKEN            = "/token"
	VALIDATE         = "/validate"
)

type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type AuthTokenProxy struct {
	clientID     string
	clientSecret string
	authToken    string
	refreshToken string
}

// NewAuthTokenProxy
func NewAuthTokenProxy(clientID string, clientSecret string, authToken string, refreshToken string) *AuthTokenProxy {
	return &AuthTokenProxy{
		clientID:     clientID,
		clientSecret: clientSecret,
		authToken:    authToken,
		refreshToken: refreshToken,
	}
}

// getAuthToken
func (a *AuthTokenProxy) GetAuthToken() (string, error) {
	if !a.validateAuthToken() {
		if err := a.refreshAuthToken(); err != nil {
			return "", err
		}
	}
	return a.authToken, nil
}

// refreshAuthToken
func (a *AuthTokenProxy) refreshAuthToken() error {
	req, err := http.NewRequest("POST", TWITCH_OAUTH_API+TOKEN, bytes.NewReader([]byte{}))
	req.Header.Add("Client-Id", a.clientID)

	q := req.URL.Query()
	q.Add("grant_type", "refresh_token")
	q.Add("client_id", a.clientID)
	q.Add("client_secret", a.clientSecret)
	q.Add("refresh_token", a.refreshToken)
	req.URL.RawQuery = q.Encode()

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Critical: Failed to refresh Token")
	}
	if response.ContentLength <= 0 {
		return errors.New("Error refreshing token. Unexpected content length")
	}

	respBody := make([]byte, response.ContentLength)
	response.Body.Read(respBody)

	newTokens := new(TokenResponse)
	if err := json.Unmarshal(respBody, newTokens); err != nil {
		return err
	}

	a.authToken = newTokens.AccessToken
	a.refreshToken = newTokens.RefreshToken

	return nil
}

// validateAuthToken
func (a *AuthTokenProxy) validateAuthToken() bool {
	req, err := http.NewRequest(http.MethodGet, TWITCH_OAUTH_API+VALIDATE, bytes.NewReader([]byte{}))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.authToken))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK || response.ContentLength <= 0 {
		return false
	}

	respBody := make([]byte, response.ContentLength)
	response.Body.Read(respBody)
	return true
}

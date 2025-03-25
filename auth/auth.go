package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GrantTypes struct {
	ClientCredentials     ClientCredentials
	ResourceOwnerPassword ResourceOwnerPassword
}

type ClientCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	Token           AccessToken
}

type ResourceOwnerPassword struct {
	Username string
	Password string
	Token    RefreshToken
}

type RefreshToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func (g ClientCredentials) ObtainAccessToken(host string) (AccessToken, error) {
	url := fmt.Sprintf("%s/api/oauth/token", host)
	payload := strings.NewReader(fmt.Sprintf("{\n  \"grant_type\": \"client_credentials\",\n  \"client_id\": \"%s\",\n  \"client_secret\": \"%s\"\n}", g.AccessKeyID, g.SecretAccessKey))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return AccessToken{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AccessToken{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return AccessToken{}, nil
	}

	var token AccessToken

	err = json.Unmarshal(body, &token)
	if err != nil {
		return AccessToken{}, err
	}
	return token, nil
}

func (g ResourceOwnerPassword) ObtainAccessToken(host string) (RefreshToken, error) {
	url := fmt.Sprintf("%s/api/oaut/token", host)
	payload := strings.NewReader(fmt.Sprintf("{\n  \"client_id\": \"administration\",\n  \"grant_type\": \"password\",\n  \"scopes\": \"write\",\n  \"username\": \"%s\",\n  \"password\": \"%s\"\n}", g.Username, g.Password))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return RefreshToken{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RefreshToken{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RefreshToken{}, err
	}

	var token RefreshToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		return RefreshToken{}, err
	}

	return token, nil
}

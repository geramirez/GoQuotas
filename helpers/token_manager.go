package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TokenRes struct {
	// Basic token struct that CF url returns
	AccessToken  string `json:"access_token"`
	Expires      int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	JTI          string `json:"jti"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	// Modified token struct with a time stamp to check if it's expired
	TokenRes
	CreatedTime int
}

type APIResponse struct {
	// Basic API struct used in CF api responses
	TotalResults int    `json:"total_results"`
	TotalPages   int    `json:"total_pages"`
	PrevUrl      string `json:"prev_url"`
	NextUrl      string `json:"next_url"`
}

func config_token_request() *http.Request {
	// Configure a new token request
	token_url := fmt.Sprintf("https://uaa.%s/oauth/token", os.Getenv("API_URL"))
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", os.Getenv("CF_USERNAME"))
	data.Set("password", os.Getenv("CF_PASSWORD"))
	req, _ := http.NewRequest("POST", token_url, bytes.NewBufferString(data.Encode()))
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authorization", "Basic Y2Y6")
	return req
}

func NewToken() *Token {
	// Initalize a new token
	var token Token
	token.get_token()
	return &token
}

func (token *Token) get_token() {
	// Get a new token
	req := config_token_request()
	client := &http.Client{}
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if json.Unmarshal(body, token) != nil {
		fmt.Println("Error")
	}
	token.CreatedTime = int(time.Now().Unix())
}

func (token *Token) update_token() {
	// Check if token has expired, if so updates the token
	if int(time.Now().Unix())-token.CreatedTime > token.Expires {
		// replace with a token refresher
		token.get_token()
	}
}

func (token *Token) make_request(req_url string) *http.Response {
	// Makes a request to the specific url with the token
	token.update_token()
	req, _ := http.NewRequest("GET", req_url, nil)
	req.Header.Set("authorization", fmt.Sprintf("bearer %s", token.AccessToken))
	client := &http.Client{}
	res, _ := client.Do(req)
	return res
}

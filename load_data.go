package main

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
	AccessToken  string `json:"access_token"`
	Expires      int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	JTI          string `json:"jti"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	TokenRes
	CreatedTime int
}

func config_request() *http.Request {
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

func update_token(req *http.Request, token *Token) {
	client := &http.Client{}
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if json.Unmarshal(body, &token) != nil {
		fmt.Println("Error")
	}
	token.CreatedTime = int(time.Now().Unix())
}

func get_token(token *Token) string {
	var action string
	if token.CreatedTime == 0 {
		req := config_request()
		update_token(req, token)
		action = "Updated"
	} else if int(time.Now().Unix())-token.CreatedTime > token.Expires {
		req := config_request()
		update_token(req, token)
		action = "Refreshed"
	} else {
		action = "None"
	}
	return action
}

func main() {
	var token Token
	fmt.Println(get_token(&token))
	fmt.Println(get_token(&token))
}

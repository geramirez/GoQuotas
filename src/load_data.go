package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
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

type QuotaMetaData struct {
	Guid    string `json:"guid"`
	Url     string `json:"url"`
	Created string `json:"created_at"`
	Updated string `json:"updated_at"`
}

type QuotaEntity struct {
	Name                    string `json:"name"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	TotalServices           int    `json:"total_services"`
	TotalRoutes             int    `json:"total_routes"`
	MemoryLimit             int    `json:"memory_limit"`
	TrialDBAllowed          bool   `json:"trial_db_allowed"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
}

type QuotaResource struct {
	MetaData QuotaMetaData `json:"metadata"`
	Entity   QuotaEntity   `json:"entity"`
}

type APIResponse struct {
	TotalResults int    `json:"total_results"`
	TotalPages   int    `json:"total_pages"`
	PrevUrl      string `json:"prev_url"`
	NextUrl      string `josn:"next_url"`
}

type QuotaAPIResponse struct {
	APIResponse
	Resources []QuotaResource `json:"resources"`
}

func config_token_request() *http.Request {
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
	if json.Unmarshal(body, token) != nil {
		fmt.Println("Error")
	}
	token.CreatedTime = int(time.Now().Unix())
}

func get_token(token *Token) string {
	var action string
	if token.CreatedTime == 0 {
		req := config_token_request()
		update_token(req, token)
		action = "Updated"
	} else if int(time.Now().Unix())-token.CreatedTime > token.Expires {
		req := config_token_request()
		update_token(req, token)
		action = "Refreshed"
	} else {
		action = "None"
	}
	return action
}

func make_request(req_url string, token *Token) *http.Response {
	get_token(token)
	req, _ := http.NewRequest("GET", req_url, nil)
	req.Header.Set("authorization", fmt.Sprintf("bearer %s", token.AccessToken))
	client := &http.Client{}
	res, _ := client.Do(req)
	return res
}

func get_quotas(token *Token) *QuotaAPIResponse {
	req_url := fmt.Sprintf("https://api.%s%s", os.Getenv("API_URL"), "/v2/quota_definitions")
	res := make_request(req_url, token)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var quotas QuotaAPIResponse
	if json.Unmarshal(body, &quotas) != nil {
		fmt.Println("Error")
	}
	return &quotas
}

func main() {
	// Initalize Token
	var token Token
	get_token(&token)
	// Get Quotas data
	quotas := get_quotas(&token)
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	for _, quota := range quotas.Resources {
		err = db.QueryRow(
			`INSERT INTO quotas(guid, name) VALUES($1, $2)`,
			quota.MetaData.Guid,
			quota.Entity.Name).Scan()
		if err != nil {
			fmt.Println(err)
		}
		err = db.QueryRow(
			"INSERT INTO quotadata(guid, memory, date) VALUES($1, $2, $3)",
			quota.MetaData.Guid,
			quota.Entity.MemoryLimit,
			time.Now().Format("2006-01-02")).Scan()
		if err != nil {
			fmt.Println(err)
		}
	}
	rows, _ := db.Query("SELECT name FROM quotas")
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Name:", name)
	}
}

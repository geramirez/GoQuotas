package helpers

import (
	"bytes"
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
	NextUrl      string `josn:"next_url"`
}

type QuotaAPIResponse struct {
	// Struct of API response for quota data
	APIResponse
	Resources []QuotaResource `json:"resources"`
}

type QuotaMetaData struct {
	// Quota meta data struct returned from the CF api
	Guid    string `json:"guid"`
	Url     string `json:"url"`
	Created string `json:"created_at"`
	Updated string `json:"updated_at"`
}

type QuotaEntity struct {
	// Quota entity sturct returned from the CF api
	Name                    string `json:"name"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	TotalServices           int    `json:"total_services"`
	TotalRoutes             int    `json:"total_routes"`
	MemoryLimit             int    `json:"memory_limit"`
	TrialDBAllowed          bool   `json:"trial_db_allowed"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
}

type QuotaResource struct {
	// Quota resource struct returned from the CF api, composed
	// composed of metadata and entity data.
	MetaData QuotaMetaData `json:"metadata"`
	Entity   QuotaEntity   `json:"entity"`
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


func (token *Token) GetQuotas() *QuotaAPIResponse {
	// Get a list of quotas and converts it to the QuotaAPIResponse struct
	req_url := fmt.Sprintf("https://api.%s%s", os.Getenv("API_URL"), "/v2/quota_definitions")
	res := token.make_request(req_url)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var quotas QuotaAPIResponse
	if json.Unmarshal(body, &quotas) != nil {
		fmt.Println("Error")
	}
	return &quotas
}

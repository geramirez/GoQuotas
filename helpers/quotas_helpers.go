package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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

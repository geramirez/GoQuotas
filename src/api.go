package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type MemoryInstance struct {
	Days int `json:"days"`
	Size int `json:"days"`
}

type Quota struct {
	Name   string           `json:"name"`
	Guid   string           `json:"guid"`
	Memory []MemoryInstance `json:"memory"`
}

type Quotas []Quota

// Open DB client

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/quotas", QuotaIndex)
	router.GET("/quotas/:guid", QuotaDetails)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func QuotaIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	quotas := Quotas{
		Quota{Name: "Cloud.gov"},
		Quota{Name: "FOIA.gov"},
	}
	json.NewEncoder(w).Encode(quotas)
}

func QuotaDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	json.NewEncoder(w).Encode(Quota{Name: "Cloud.gov"})

}

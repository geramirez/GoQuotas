package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strings"
)

type memory struct {
	MemoryLimit int `json:"memory"`
	Days        int `json:"days"`
}

type Quota struct {
	Guid   string   `json:"guid"`
	Name   string   `json:"name"`
	Memory []memory `json:"data"`
}

type Quotas []Quota

func main() {

	router := httprouter.New()
	//TODO: Set to be main route
	//TODO: Add authentication
	router.GET("/", Index)
	router.GET("/api/quotas", QuotaIndex)
	router.GET("/api/quotas/:guid", QuotaDetails)
	router.ServeFiles("/static/*filepath", http.Dir("static/"))
	router.ServeFiles("/dist/*filepath", http.Dir("static/dist/"))

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func QuotaIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//TODO: Convert to function
	query_values := r.URL.Query()
	since := query_values.Get("since")
	if since == "" {
		since = "1970-01-01"
	}
	until := query_values.Get("until")
	if until == "" {
		until = "2050-01-01"
	}
	//TODO: Setup to only open once
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var quotas string
	err = db.QueryRow(`
		with quota_details_agg as (
		  select * from get_quotas_details($1, $2)
		)
		SELECT json_agg(t) AS elements FROM (
		  SELECT guid, name, cost, data
		  FROM quota_details_agg
		) t`,
		string(since),
		string(until),
	).Scan(&quotas)
	var data string
	if err == nil {
		data = fmt.Sprintf(`{"Quotas": %s}`, quotas)
	} else {
		fmt.Println(err)
		data = `{"Quotas": []}`
	}
	fmt.Fprint(w, data)
}

func QuotaDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//TODO: Convert to function
	query_values := r.URL.Query()
	since := query_values.Get("since")
	if since == "" {
		since = "1970-01-01"
	}
	until := query_values.Get("until")
	if until == "" {
		until = "2050-01-01"
	}

	//TODO: Setup to only open once
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var quotas string
	err = db.QueryRow(`
		SELECT details
		FROM get_quotas($1, $2)
		WHERE guid = $3`,
		since,
		until,
		ps.ByName("guid"),
	).Scan(&quotas)
	fmt.Fprint(w, strings.TrimRight(strings.TrimLeft(quotas, "["), "]"))
}

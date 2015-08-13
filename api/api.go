package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strings"
)

// Quota-related structs
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

// Context functions
type app_context struct {
	db *sql.DB
}

// Query function
func get_dates(r *http.Request) (string, string) {
	query_values := r.URL.Query()
	since := query_values.Get("since")
	if since == "" {
		since = "1970-01-01"
	}
	until := query_values.Get("until")
	if until == "" {
		until = "2050-01-01"
	}
	return since, until
}

// Route functions
func (app *app_context) QuotaList(w http.ResponseWriter, r *http.Request) {
	// Route function for a quota list endpoint
	since, until := get_dates(r)
	var quotas string
	app.db.QueryRow(`
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
	if quotas == "" {
		quotas = "[]"
	}
	fmt.Fprint(w, fmt.Sprintf(`{"Quotas": %s}`, quotas))
}

func (app *app_context) QuotaDetails(w http.ResponseWriter, r *http.Request) {
	// Route function for quota details endpoint
	guid := mux.Vars(r)["guid"]
	since, until := get_dates(r)

	var quotas string
	app.db.QueryRow(`
		SELECT details
		FROM get_quotas($1, $2)
		WHERE guid = $3`,
		since,
		until,
		guid,
	).Scan(&quotas)
	fmt.Println(quotas)
	if quotas == "" {
		quotas = "{}"
	}
	fmt.Fprint(w, strings.TrimRight(strings.TrimLeft(quotas, "["), "]"))
}

// Route functions
func (app *app_context) CSVView(w http.ResponseWriter, r *http.Request) {
	// Route function for a quota list endpoint
	since, until := get_dates(r)
	var csv string
	app.db.QueryRow(`
		COPY (SELECT guid, name, cost FROM get_quotas_details($1, $2)) TO STDOUT WITH CSV HEADER;
		`,
		string(since),
		string(until),
	).Scan(&csv)
	fmt.Println(csv)
	fmt.Fprint(w, fmt.Sprintf(`%s`, csv))
}

func main() {

	//Open Database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	context := &app_context{db: db}
	router := mux.NewRouter()
	//TODO: Add authentication
	router.HandleFunc("/api/quotas", context.QuotaList)
	router.HandleFunc("/api/quotas/{guid}", context.QuotaDetails)
	router.HandleFunc("/quotas.csv", context.CSVView)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))
	log.Fatal(http.ListenAndServe(":8080", router))
}

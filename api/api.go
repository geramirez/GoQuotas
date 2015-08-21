package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

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
	if quotas == "" {
		quotas = "{}"
	}
	fmt.Fprint(w, strings.TrimRight(strings.TrimLeft(quotas, "["), "]"))
}

// Route functions
func (app *app_context) CSVView(w http.ResponseWriter, r *http.Request) {
	// Route function for a quota list endpoint
	since, until := get_dates(r)
	rows, _ := app.db.Query(`
		select guid, name, cost from get_quotas_details($1, $2);
		`,
		since,
		until)
	defer rows.Close()
	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	wr.Write([]string{"guid", "name", "cost"})
	for rows.Next() {
		var guid, name string
		var cost string
		rows.Scan(&guid, &name, &cost)
		wr.Write([]string{guid, name, cost})
	}
	wr.Flush()
	w.Header().Set("Content-Type", "text/csv")
	fmt.Fprint(w, string(b.Bytes()))
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

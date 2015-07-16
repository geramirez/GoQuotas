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
	router.GET("/", Index)
	router.GET("/api/quotas", QuotaIndex)
	router.GET("/api/quotas/:guid", QuotaDetails)
	router.ServeFiles("/static/*filepath", http.Dir("static/"))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func QuotaIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var quotas string
	err = db.QueryRow(`
		SELECT json_agg(t) AS elements FROM (SELECT guid, name, data FROM quotas_view ) t
	`).Scan(&quotas)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, quotas)
}

func QuotaDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var quotas string
	err = db.QueryRow(
		"SELECT details FROM quota_details WHERE guid = $1",
		ps.ByName("guid"),
	).Scan(&quotas)
	fmt.Fprint(w, strings.TrimRight(strings.TrimLeft(quotas, "["), "]"))
}

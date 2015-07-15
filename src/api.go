package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Quota struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

type Quotas []Quota

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

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	rows, _ := db.Query("SELECT * FROM quotas")
	var quotas Quotas
	for rows.Next() {
		var name string
		var guid string
		err = rows.Scan(&guid, &name)
		if err != nil {
			fmt.Println(err)
		}
		quotas = append(quotas, Quota{guid, name})
	}
	json.NewEncoder(w).Encode(quotas)
}

func QuotaDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	rows, _ := db.Query("SELECT * FROM quotas WHERE guid = $1", ps.ByName("guid"))
	var quota Quota
	for rows.Next() {
		var name string
		var guid string
		err = rows.Scan(&guid, &name)
		if err != nil {
			fmt.Println(err)
		}
		quota = Quota{guid, name}
	}
	json.NewEncoder(w).Encode(quota)

}

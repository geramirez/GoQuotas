package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"

	helpers "github.com/ramirezg/GoQuotas/helpers"
)

func main() {
	// Collect quotas and quota mememory data
	token := helpers.NewToken()
	quotas := token.GetQuotas()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	for _, quota := range quotas.Resources {
		err = db.QueryRow(`
      INSERT INTO quotas(guid, name) VALUES($1, $2)`,
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
}

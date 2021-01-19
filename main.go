package main

import (
	"database/sql"
	"log"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"os"
	"net/http"
)
func main() {
	http.HandleFunc("/", handler)
	log.Println("Listening on localhost:8282")
	log.Fatal(http.ListenAndServe(":8282", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, getStatus())
}

func getStatus() (status string) {

		// Set status "down" as default
		status = "down"

		//connect to the database
		host := os.Getenv("SQLSERVER_HOST")
		log.Println("SQL Server host:", host)
		conn := fmt.Sprintf("sqlserver://%s/?database=master&connection+timeout=30",host)
		db, err := sql.Open("sqlserver", conn)

		// Set status "down" for db down
		if err != nil {
			log.Println(err)
			status = "down"	
		}  else {
			// Query DB for role
			query := "select top 1 mirroring_role_desc from sys.database_mirroring where mirroring_role is not null"
			rows, err := db.Query(query)
			if err != nil {
				log.Println(err)
			}
			defer rows.Close()

			// Read return result
			var strrow string
			for rows.Next() {
				err = rows.Scan(&strrow)
			}

			// Set status as none
			if strrow == "" {
				status = "none"
			} else if strrow == "PRINCIPAL" {
				status = "principal"
			} else if strrow == "MIRROR" {
				status = "mirror"
			}
		}

		log.Println("I am: " + status)
		return status
}
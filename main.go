package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	//"errors"
)

func main() {

	// checkEnvironment()

	http.HandleFunc("/", handler)
	log.Println("Listening on localhost:8282")
	log.Fatal(http.ListenAndServe(":8282", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Check SQL Role and set 200 if principal
	status := getStatus()
	if status != "principal" {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Return status
	fmt.Fprintf(w, status)
}

func getStatus() (status string) {

	// Set status "down" as default
	status = "down"

	//connect to the database
	host := os.Getenv("SQLSERVER_HOST")
	log.Println("SQL Server host:", host)
	conn := "sqlserver://" + host + "/?database=master&connection+timeout=30"
	db, err := sql.Open("sqlserver", conn)

	// Set status "down" for db down
	if err != nil {
		log.Println(err)
		status = "down"
	} else {
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

// Ensures all env var for running this are defined
func checkEnvironment() {

	envvars := []string{
		"SQLSERVER_HOST",
		"SQLSERVER_PORT",
		"LISTEN_PORT",
	}

	checkEnvironmentVariable(envvars)
}

// Checks is env va exists, else logs fatal
func checkEnvironmentVariable(envvars []string) {
	for _, envvar := range envvars {
		if os.Getenv(envvar) == "" {
			log.Fatal("Error, environent vaiable " + envvar + " was not defined. Exitting.")
		}
	}
}

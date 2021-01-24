package core

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Config contains the Server Configurations required for the server start
type Config struct {
	// SQLServerHost is the IP address of the SQL Server to Monitor
	SQLServerHost string

	// SQLServerPort is the connect port of the SQL Server to Monitor
	SQLServerPort uint16

	// ListenPort is the web server port this server should listen on. Default: Port 8282
	ListenPort uint16
}

func GetServiceConfig() (cfg *Config) {
	cfg = &Config{
		SQLServerHost: "127.0.0.1",
		SQLServerPort: 1433,
		ListenPort:    8282,
	}

	if sqlHost := os.Getenv("SQLSERVER_HOST"); sqlHost != "" {
		cfg.SQLServerHost = os.Getenv("SQLSERVER_HOST")
	}

	if sqlPort := os.Getenv("SQLSERVER_PORT"); sqlPort != "" {
		intPort, err := strconv.Atoi(sqlPort)
		if err != nil {
			log.Fatalf("Provided value for SQLSERVER_PORT, %s, cannot be converted to an integer", sqlPort)
		}
		cfg.SQLServerPort = uint16(intPort)
	}

	if listenPort := os.Getenv("LISTEN_PORT"); listenPort != "" {
		intPort, err := strconv.Atoi(listenPort)
		if err != nil {
			log.Fatalf("Provided value for LISTEN_PORT, %s, cannot be converted to an integer", listenPort)
		}
		cfg.ListenPort = uint16(intPort)
	}

	return cfg
}

func (cfg *Config) Handler(w http.ResponseWriter, r *http.Request) {

	// Check SQL Role and set 200 if principal
	mirroring, err := cfg.getMirroringStatus()
	if err != nil {
		log.Println(err)
		http.Error(w, JsonErrorIt(err.Error()), http.StatusInternalServerError)
	} else {
		if mirroring.OverallMirroringRole != "principal" {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Return status
		fmt.Fprintf(w, JsonPrintIt(&JsonResponse{Data: mirroring}))
	}
}

func (cfg *Config) getMirroringStatus() (mirroring *JsonMirroring, err error) {

	// Initialise response
	mirroring = &JsonMirroring{
		OverallMirroringRole: "none",
		DatabasesMirroring:   make([]*JsonDbMirroring, 0),
	}

	// Connect to the database
	conn := fmt.Sprintf("sqlserver://%s:%d/?database=master&connection+timeout=30", cfg.SQLServerHost, cfg.SQLServerPort)
	db, err := sql.Open("sqlserver", conn)
	defer db.Close()

	// Return on error
	if err != nil {
		return mirroring, err
	}

	// Query DB for db names roles
	query := "SELECT m.name, d.database_id, d.mirroring_role_desc FROM sys.database_mirroring d, sys.databases m WHERE mirroring_role is not null and m.database_id = d.database_id"
	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if rows != nil {
		mirrorRolesMap := make(map[string]struct{}, 0)

		for rows.Next() {

			var dbName string
			var dbID int
			var dbRole string

			// Read data from row
			err = rows.Scan(&dbName, &dbID, &dbRole)
			if err != nil {
				return nil, err
			}

			dbRole = strings.ToLower(dbRole)
			dbMirroring := &JsonDbMirroring{
				Name:          dbName,
				MirroringRole: dbRole,
			}

			//
			if mirroring.OverallMirroringRole == "none" {
				// store the role of the first db in the overall
				mirroring.OverallMirroringRole = dbRole
			} else if mirroring.OverallMirroringRole == "principal" && dbRole != "principal" {
				mirroring.OverallMirroringRole = dbRole
			}
			mirroring.DatabasesMirroring = append(mirroring.DatabasesMirroring, dbMirroring)
			mirrorRolesMap[dbRole] = struct{}{}

		}

		// multiple db mirroring roles result in a "mixed" overall mirroring role
		if len(mirrorRolesMap) > 1 {
			mirroring.OverallMirroringRole = "mixed"
		}
	}

	log.Println("I am: " + mirroring.OverallMirroringRole)
	return mirroring, nil
}

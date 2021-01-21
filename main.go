package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	cfg := getServiceConfig()
	http.HandleFunc("/", cfg.handler)
	listenOn := fmt.Sprintf("0.0.0.0:%d", cfg.ListenPort)
	log.Printf("Listening on %s", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}

func (cfg *Config) handler(w http.ResponseWriter, r *http.Request) {

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
		fmt.Fprintf(w, PrettyPrintJson(mirroring))
	}
}

func (cfg *Config) getMirroringStatus() (mirroring *JsonMirroring, err error) {
	mirroring = &JsonMirroring{
		OverallMirroringRole: "none",
		DatabasesMirroring:   make([]*JsonDbMirroring, 0),
	}
	//connect to the database
	conn := fmt.Sprintf("sqlserver://%s:%d/?database=master&connection+timeout=30", cfg.SqlServerHost, cfg.SqlServerPort)
	db, err := sql.Open("sqlserver", conn)

	if err != nil {
		return mirroring, err
	} else {
		// Query DB for db names roles
		query := "SELECT m.name, d.database_id, d.mirroring_role_desc FROM sys.database_mirroring d, sys.databases m WHERE mirroring_role is not null and m.database_id = d.database_id"
		rows, err := db.Query(query)
		if err != nil {
			return nil, err
		}

		if rows != nil {
			defer rows.Close()
			mirrorRolesMap := make(map[string]struct{}, 0)
			for rows.Next() {
				var dbName string
				var dbId int
				var dbRole string
				err = rows.Scan(&dbName, &dbId, &dbRole)
				if err != nil {
					return nil, err
				}
				dbRole = strings.ToLower(dbRole)
				dbMirroring := &JsonDbMirroring{
					Name:          dbName,
					MirroringRole: dbRole,
				}
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
	}
	log.Println("I am: " + mirroring.OverallMirroringRole)
	return mirroring, nil
}

func PrettyPrintJson(v interface{}) string {
	outBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("ERROR: PrettyPrintJson error '%s' for v='%+v'", err.Error(), v)
	}
	outBytes = append(outBytes, "\n"...)
	return fmt.Sprintf("%-512s", string(outBytes)) // Pad right to avoid "Friendly HTTP error messages" issue in IE & Chrome
}

func JsonErrorIt(msg string) string {
	jsonErr := JsonResponse{
		Error: msg,
	}
	outBytes, err := json.MarshalIndent(jsonErr, "", "  ")
	if err != nil {
		log.Printf("ERROR: JsonErrorIt marshal error '%s'", err.Error())
	}
	return fmt.Sprintf("%-512s", outBytes) // Pad right to avoid "Friendly HTTP error messages" issue in IE & Chrome
}

func getServiceConfig() (cfg *Config) {
	cfg = &Config{
		SqlServerHost: "127.0.0.1",
		SqlServerPort: 1433,
		ListenPort:    8282,
	}

	if sqlHost := os.Getenv("SQLSERVER_HOST"); sqlHost != "" {
		cfg.SqlServerHost = os.Getenv("SQLSERVER_HOST")
	}

	if sqlPort := os.Getenv("SQLSERVER_PORT"); sqlPort != "" {
		intPort, err := strconv.Atoi(sqlPort)
		if err != nil {
			log.Fatalf("Provided value for SQLSERVER_PORT, %s, cannot be converted to an integer", sqlPort)
		}
		cfg.SqlServerPort = uint16(intPort)
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

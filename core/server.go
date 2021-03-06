package core

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/sys/windows/svc/eventlog"
)

// VERSION is the version of the server
var VERSION string = "0.0.0"
var elog *eventlog.Log

type Runnable interface {
	Start() error
	Stop() error
}

// Server wrapps the check mirror server in a server object that implements
type Server struct {
	config Config
}

// NewServer creates and returns a new server
func NewServer() *Server {
	var err error
	elog, err = eventlog.Open("CheckMirror")
	if err != nil {
		log.Fatalf("Error initializing elog: %s", err.Error())
	}
	if elog == nil {
		log.Fatalf("Error initializing elog, can't be nil!")
	}
	myConfig := GetServiceConfig()

	server := &Server{
		config: *myConfig,
	}

	return server
}

// Start stars the server
func (s *Server) Start() {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.ListenPort)
	elog.Info(1, fmt.Sprintf("Server (%s) started on %s", VERSION, addr))
	log.Printf("Server (%s) started on %s", VERSION, addr)
	server := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(s.defaultHandler),
	}

	err := server.ListenAndServe()
	if err != nil {
		elog.Error(100, fmt.Sprintf("ListenAndServe returned error: %s", err.Error()))
		log.Fatalf("ListenAndServe returned error: %s", err)
	}
}

// Stop Stops the server
func (*Server) Stop() error {
	elog.Close()
	// Stop the service here
	return nil
}

// defaultHandler The default handler for the service. Returns the mirror status
func (s *Server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	//Check SQL Role and set 200 if principal
	mirroring, err := s.getMirroringStatus()
	if err != nil {
		log.Println(fmt.Sprintf("Unable to get mirroring status, error: %s", err.Error()))
		elog.Error(101, fmt.Sprintf("Unable to get mirroring status, error: %s", err.Error()))
		http.Error(w, JsonErrorIt(err.Error()), http.StatusInternalServerError)
	} else {
		if mirroring.OverallMirroringRole != "principal" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// Return status
		fmt.Fprintf(w, JsonPrintIt(&JsonResponse{Data: mirroring}))
	}
}

func (s *Server) getMirroringStatus() (mirroring *JsonMirroring, err error) {
	// Initialise response
	mirroring = &JsonMirroring{
		OverallMirroringRole: "none",
		DatabasesMirroring:   make([]*JsonDbMirroring, 0),
	}

	// Connect to the database
	conn := fmt.Sprintf("sqlserver://%s:%d/?database=master&connection+timeout=30", s.config.SQLServerHost, s.config.SQLServerPort)
	// log.Printf("The connection string is: %s", conn)
	db, err := sql.Open("sqlserver", conn)

	if db.Ping() != nil {
		return mirroring, fmt.Errorf("Could not connect to database on %s:%d", s.config.SQLServerHost, s.config.SQLServerPort)
	}

	defer db.Close()
	// Query DB for db names roles
	query := "SELECT m.name, d.database_id, d.mirroring_role_desc FROM sys.database_mirroring d, sys.databases m WHERE mirroring_role is not null and m.database_id = d.database_id"
	rows, err := db.Query(query)

	if err != nil {
		return mirroring, err
	}

	if rows != nil {
		defer rows.Close()

		// Iterate through result set to construct json response
		mirrorRolesMap := make(map[string]struct{}, 0)
		for rows.Next() {
			// Read data from row
			var dbName string
			var dbID int
			var dbRole string
			err = rows.Scan(&dbName, &dbID, &dbRole)
			if err != nil {
				return mirroring, err
			}

			dbRole = strings.ToLower(dbRole)
			dbMirroring := &JsonDbMirroring{
				Name:          dbName,
				MirroringRole: dbRole,
			}

			// Get Overall Mirroring Role
			if mirroring.OverallMirroringRole == "none" {
				// store the role of the first db in the overall
				mirroring.OverallMirroringRole = dbRole
			} else if mirroring.OverallMirroringRole == "principal" && dbRole != "principal" {
				mirroring.OverallMirroringRole = dbRole
			}

			// Append database and role to list
			mirroring.DatabasesMirroring = append(mirroring.DatabasesMirroring, dbMirroring)
			mirrorRolesMap[dbRole] = struct{}{}

		}

		// multiple db mirroring roles result in a "mixed" overall mirroring role
		if len(mirrorRolesMap) > 1 {
			mirroring.OverallMirroringRole = "mixed"
		}
	}

	elog.Info(2, fmt.Sprintf("Got '%s' mirroring status from %s:%d\n", mirroring.OverallMirroringRole, s.config.SQLServerHost, s.config.SQLServerPort))
	log.Printf("Got '%s' mirroring status from %s:%d\n", mirroring.OverallMirroringRole, s.config.SQLServerHost, s.config.SQLServerPort)
	return mirroring, nil
}

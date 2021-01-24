package core

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// VERSION is the version of the server
var VERSION string = "0.0.0"

// Service interface defines methods that need to be implemented for running this app as a windows service
type Service interface {
	start() bool
	stop() bool
	status() string
}

// Server wrapps the check mirror server in a server object that implements
type Server struct {
	config Config
}

// NewServer creates and returns a new server
func NewServer() *Server {

	myConfig := GetServiceConfig()

	server := &Server{
		config: *myConfig,
	}

	return server
}

// Start stars the server
func (s *Server) Start() {

	addr := fmt.Sprintf("0.0.0.0:%d", s.config.ListenPort)
	server := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(s.defaultHandler),
	}

	log.Fatal(server.ListenAndServe())
}

// defaultHandler The default handler for the service. Returns the mirror status
func (s *Server) defaultHandler(w http.ResponseWriter, r *http.Request) {

	//Check SQL Role and set 200 if principal
	mirroring, err := s.getMirroringStatus()
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

func (s *Server) getMirroringStatus() (mirroring *JsonMirroring, err error) {

	log.Println("Getting Mirroring Status:")
	// Initialise response
	mirroring = &JsonMirroring{
		OverallMirroringRole: "none",
		DatabasesMirroring:   make([]*JsonDbMirroring, 0),
	}

	// Connect to the database
	conn := fmt.Sprintf("sqlserver://%s:%d/?database=master&connection+timeout=30", s.config.SQLServerHost, s.config.SQLServerPort)
	log.Printf("The connection string is: %s", conn)
	db, err := sql.Open("sqlserver", conn)
	defer db.Close()

	if db.Ping() != nil {
		return mirroring, fmt.Errorf("Could not connect to database on %s:%d", s.config.SQLServerHost, s.config.SQLServerPort)
	}

	// Query DB for db names roles
	query := "SELECT m.name, d.database_id, d.mirroring_role_desc FROM sys.database_mirroring d, sys.databases m WHERE mirroring_role is not null and m.database_id = d.database_id"
	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		return mirroring, err
	}

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

	log.Println("I am: " + mirroring.OverallMirroringRole)
	return mirroring, nil
}

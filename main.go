package main

import (
	// "database/sql"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"writeameer/checkmirror/core"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	cfg := core.GetServiceConfig()
	http.HandleFunc("/", cfg.Handler)
	listenOn := fmt.Sprintf("0.0.0.0:%d", cfg.ListenPort)
	log.Printf("Listening on %s", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}

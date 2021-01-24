package main

import (
	"writeameer/checkmirror/core"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {

	// Create Server
	checkmirror := core.NewServer()

	// Start Server
	checkmirror.Start()

}

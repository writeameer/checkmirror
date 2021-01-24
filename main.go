package main

import (
	"fmt"
	"writeameer/checkmirror/core"

	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/hlandau/easyconfig.v1"
	"gopkg.in/hlandau/service.v2"
)

const serviceName = "CheckMirror"

func main() {
	// deals with the --service.do=* options
	easyconfig.ParseFatal(nil, nil)

	service.Main(&service.Info{
		Title:       "Check SQL Mirroring Service",
		Name:        serviceName,
		Description: fmt.Sprintf("%s queries the local MS SQL Server to return an HTTP 200 status code if it is a principal", serviceName),

		RunFunc: func(smgr service.Manager) error {
			// Create & start server
			checkmirror := core.NewServer()
			go checkmirror.Start()

			// Once initialization requiring root is done, call this.
			err := smgr.DropPrivileges()
			if err != nil {
				return err
			}

			// When it is ready to serve requests, call this.
			// You must call DropPrivileges first.
			smgr.SetStarted()

			// Optionally set a status.
			smgr.SetStatus(fmt.Sprintf("%s: running ok", serviceName))

			// Wait until stop is requested.
			<-smgr.StopChan()

			// Do any necessary teardown.
			// ...

			// Done.
			return nil
		},
	})

}

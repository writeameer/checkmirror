package main

import (
	"writeameer/checkmirror/core"

	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/hlandau/easyconfig.v1"
	"gopkg.in/hlandau/service.v2"
)

func main() {

	easyconfig.ParseFatal(nil, nil)

	service.Main(&service.Info{
		Title:       "Checkmirror Server",
		Name:        "Checkmirror",
		Description: "Checkmirror queries the local MS SQL Server to return an HTTP 200 status code if it is a principal",

		RunFunc: func(smgr service.Manager) error {
			// Create & start server
			checkmirror := core.NewServer()
			checkmirror.Start()

			// Once initialization requiring root is done, call this.
			err := smgr.DropPrivileges()
			if err != nil {
				return err
			}

			// When it is ready to serve requests, call this.
			// You must call DropPrivileges first.
			smgr.SetStarted()

			// Optionally set a status.
			smgr.SetStatus("Checkmirror: running ok")

			// Wait until stop is requested.
			<-smgr.StopChan()

			// Do any necessary teardown.
			// ...

			// Done.
			return nil
		},
	})

}

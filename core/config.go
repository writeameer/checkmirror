package core

import (
	"log"
	"os"
	"strconv"
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

// GetServiceConfig Returns the config for service based on env input or defaults.
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

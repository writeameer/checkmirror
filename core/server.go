package core

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

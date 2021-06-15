// Package main invokes the server package to start and run the server.
package main

import (
	"ProjectGoLive/server"
)

// init calls server.InitServer to initialise the variables. This will only be called once
// in the duration of the application.
func init() {
	server.InitServer()
}

// main calls server.StartServer() to run the application.
func main() {
	server.StartServer()
}

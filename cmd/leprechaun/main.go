package main

import (
	"fmt"
	"runtime"
)

// VERSION of application
const (
	VERSION = "1.1.0"
	RELEASE = "Calimero"
)

func main() {
	fmt.Println(runtime.Version())
	// daemon.Srv.AddService(&client.Client{Name: "scheduler"})
	// daemon.Srv.AddService(&server.Server{Name: "server"})
	// daemon.Srv.AddService(&cron.Cron{Name: "cron"})
	// daemon.Srv.Run()
}

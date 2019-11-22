// +build !remote

package main

import (
	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/cron"
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/server"
)

func main() {
	daemon.Srv.AddService(&client.Client{Name: "scheduler"})
	daemon.Srv.AddService(&server.Server{Name: "server"})
	daemon.Srv.AddService(&cron.Cron{Name: "cron"})
	daemon.Srv.Run()
}

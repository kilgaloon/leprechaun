// +build !remote

package main

import (
	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/cron"
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/server"
)

func main() {
	daemon := daemon.Init()

	if daemon != nil {
		daemon.AddService(&client.Client{Name: "scheduler"})
		daemon.AddService(&server.Server{Name: "server"})
		daemon.AddService(&cron.Cron{Name: "cron"})

		daemon.Run(nil)
	}
}

// +build remote

package main

import (
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/remote"
)

func main() {
	daemon.Srv.AddService(&remote.Remote{Name: "remote"})
	daemon.Srv.Run()
}

// +build remote

package main

import (
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/remote"
)

func main() {
	daemon := daemon.Init()
	
	if daemon != nil {
		daemon.AddService(&remote.Remote{Name: "remote"})
		daemon.Run()
	}
	
}

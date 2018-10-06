package main

import (
	"fmt"
	"os"

	"github.com/kilgaloon/leprechaun/api"
)

func help(r api.Registrator) {
	cmds := r.RegisterCommands()

	fmt.Println("---------- LIST OF AVAILABLE COMMANDS -----------")

	for name, cmd := range cmds {
		fmt.Println(name + " - " + cmd.String())
	}

	os.Exit(0)
}

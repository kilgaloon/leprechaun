package main

import (
	"fmt"
	"strings"

	"github.com/kilgaloon/leprechaun/api"
)

func help(r api.Registrator) {
	cmds := r.RegisterCommands()

	fmt.Println("---------- COMMANDS FOR " + strings.ToUpper(r.GetName()) + " -----------")

	for name, cmd := range cmds {
		formated := cmd.String()
		fmt.Println(name + " - " + strings.Replace(formated, "{agent}", r.GetName(), -1))
	}

	fmt.Println("-------------------------------------------------")
}

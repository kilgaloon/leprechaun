package main

import (
	"fmt"
	"strings"

	"github.com/kilgaloon/leprechaun/api"
)

func help(r api.Registrator) {
	cmds := r.RegisterCommands()

	fmt.Println("---------- COMMANDS FOR " + strings.ToUpper(r.GetName()) + " -----------")

	fmt.Println(r.GetName() + ":start - Start agent")
	fmt.Println(r.GetName() + ":stop - Stop agent")

	for name, cmd := range cmds {
		formated := cmd.String()
		fmt.Println(name + " - " + strings.Replace(formated, "{agent}", r.GetName(), -1))
	}

	fmt.Println("-------------------------------------------------")
}

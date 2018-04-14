package main

import (
	"./src/client"
	"flag"
)

// VERSION of application
const (
	VERSION = "0.2.0"
	RELEASE = "Calimero"
)

func main() {

	iniPath := flag.String("ini_path", "./configs/client.ini", "Path to client .ini configuration")
	flag.Parse()

	client.Start(iniPath)
}

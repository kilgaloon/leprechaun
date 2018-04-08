package main

import (
	"flag"
	"./src/client"
)

// VERSION of application
const (
	VERSION = "0.1.0"
	RELEASE = "Calimero"
)

func main() {

	iniPath := flag.String("ini_path", "./configs/client.ini", "Path to client .ini configuration")
	flag.Parse()

	client.Start(iniPath)
}
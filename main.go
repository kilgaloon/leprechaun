package main

import (
	"flag"

	"github.com/kilgaloon/leprechaun/client"
)

// VERSION of application
const (
	VERSION = "0.6.0"
	RELEASE = "Calimero"
)

func main() {

	iniPath := flag.String("ini_path", "configs/client.ini", "Path to client .ini configuration")
	flag.Parse()

	client.Start(iniPath)
}

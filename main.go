package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/server"
)

// VERSION of application
const (
	VERSION = "0.6.0"
	RELEASE = "Calimero"
)

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	shutdownSignal := make(chan os.Signal, 1)

	clientIniPath := flag.String("client_ini_path", "/etc/leprechaun/configs/client.ini", "Path to client .ini configuration")
	serverIniPath := flag.String("server_ini_path", "/etc/leprechaun/configs/server.ini", "Path to server .ini configuration")
	flag.Parse()

	client.CreateAgent(clientIniPath)
	server.CreateAgent(serverIniPath)

	switch command {
	case "client:stop":
		shutdownSignal <- client.Agent.Stop()
	default:
		go client.Agent.Start()
		go server.Agent.Start()
	}

	signal.Notify(shutdownSignal,
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGSTOP,
		syscall.SIGTERM)

	<-shutdownSignal

	os.Exit(0)
}

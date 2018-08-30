package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kilgaloon/leprechaun/client"
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

	iniPath := flag.String("ini_path", "configs/client.ini", "Path to client .ini configuration")
	flag.Parse()

	client.CreateAgent(iniPath)

	switch command {
	case "stop":
		shutdownSignal <- client.Agent.Stop()
	default:
		client.Agent.Start()
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

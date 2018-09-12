package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/config"
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

	iniPath := flag.String("ini_path", "/etc/leprechaun/configs/config.ini", "Path to .ini configuration")
	flag.Parse()

	cfg := config.BuildConfig(*iniPath)

	client.CreateAgent(cfg.GetClientConfig())
	server.CreateAgent(cfg.GetServerConfig())

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

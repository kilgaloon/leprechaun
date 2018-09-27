package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/server"
	"github.com/kilgaloon/leprechaun/socket"
)

// VERSION of application
const (
	VERSION = "0.6.0"
	RELEASE = "Calimero"
)

func main() {
	shutdownSignal := make(chan os.Signal, 1)

	iniPath := flag.String("ini_path", "/etc/leprechaun/configs/config.ini", "Path to .ini configuration")
	cmd := flag.String("cmd", "client:start", "Command for app to run")
	flag.Parse()

	cfg := config.BuildConfig(*iniPath)

	client.CreateAgent(cfg.GetClientConfig())
	server.CreateAgent(cfg.GetServerConfig())

	switch strings.Fields(*cmd)[0] {
	case "client:stop":
		shutdownSignal <- client.Agent.Stop()
	case "client:start":
		go client.Agent.Start()
	case "server:start":
		go server.Agent.Start()
	case "client":
		sock := socket.BuildSocket("var/run/client.sock")
		sock.Command(*cmd)
		os.Exit(0)
	default:
		os.Exit(0)
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

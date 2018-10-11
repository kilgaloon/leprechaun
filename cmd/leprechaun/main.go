package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/cron"
	"github.com/kilgaloon/leprechaun/server"
)

// VERSION of application
const (
	VERSION = "1.0.0-alpha"
	RELEASE = "Calimero"
)

func main() {
	shutdownSignal := make(chan os.Signal, 1)

	iniPath := flag.String("ini_path", "/etc/leprechaun/configs/config.ini", "Path to .ini configuration")
	cmd := flag.String("cmd", "run", "Command for app to run")
	flag.Parse()

	configs := config.NewConfigs()
	client.New("client", configs.New("client", *iniPath))
	server.New("server", configs.New("server", *iniPath))
	cron.New("cron", configs.New("cron", *iniPath))

	// basic leprechaun help
	if len(os.Args) > 1 {
		if os.Args[1] == "help" {
			help(client.Agent)
			help(server.Agent)
			help(cron.Agent)

			os.Exit(0)
		}
	}

	switch strings.Fields(*cmd)[0] {
	case "run":
		go client.Agent.Start()
		go server.Agent.Start()
		go cron.Agent.Start()
	case "client:stop":
		shutdownSignal <- client.Agent.Stop()
	case "client:start":
		go client.Agent.Start()
	case "server:start":
		go server.Agent.Start()
	case "cron:start":
		go cron.Agent.Start()
	case "client":
		sock := api.New(configs.GetConfig("client").GetCommandSocket())
		sock.Command(*cmd)
		os.Exit(0)
	case "server:stop":
		*cmd = "server stop"
		fallthrough
	case "server":
		sock := api.New(configs.GetConfig("server").GetCommandSocket())
		sock.Command(*cmd)
		os.Exit(0)
	case "cron:stop":
		*cmd = "cron stop"
		fallthrough
	case "cron":
		sock := api.New(configs.GetConfig("server").GetCommandSocket())
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

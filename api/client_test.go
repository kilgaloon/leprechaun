package api

import (
	"testing"

	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/stretchr/testify/assert"
)

var (
	iniFile    = "../tests/configs/config_regular.ini"
	path       = &iniFile
	cfgWrap    = config.NewConfigs()
	fakeClient = client.New("test", cfgWrap.New("test", *path), false)
)

func TestMain(t *testing.T) {
	assert.Equal(t, IsAPIRunning(), false)
	go New("").Register(fakeClient).Start()
	assert.Equal(t, IsAPIRunning(), true)
}
func TestCmd(t *testing.T) {
	c := Cmd("agent command arg")

	assert.Equal(t, c.agent(), "agent")
	assert.Equal(t, c.command(), "command")
	assert.Equal(t, c.args()[0], "arg")
}

func TestResolver(t *testing.T) {
	Resolver(Cmd("test not_exist"))
	Resolver(Cmd("test info"))
	Resolver(Cmd("test workers:list"))
	Resolver(Cmd("test workers:kill job"))

	Resolver(Cmd("agent workers:kill job"))
}

func TestRevealEndpoint(t *testing.T) {
	assert.Equal(t, "http://localhost:11401/agent/command", revealEndpoint("/{agent}/command", Cmd("agent command")))
}

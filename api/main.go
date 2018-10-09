package api

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"syscall"

	"github.com/olekukonko/tablewriter"
)

// Socket defines socket on which we listen for commands
type Socket struct {
	unixSock  string
	commands  map[string]Command
	readyChan chan bool
}

// Closure is command that will be called to execute
type Closure func(r io.Writer, args ...string) ([][]string, error)

// Command is definition of basic command that can be called with --cmd
type Command struct {
	Definition
	Closure Closure
}

// Definition of command
type Definition struct {
	Text  string
	Usage string
}

func (c Command) String() string {
	return c.Definition.Text + " - usage: " + c.Definition.Usage
}

// Registrator is agent that will be registered
type Registrator interface {
	RegisterCommands() map[string]Command
}

func (s *Socket) ready() {
	s.readyChan <- true
}

// GetCommands return array of commands
func (s *Socket) GetCommands() map[string]Command {
	return s.commands
}

// Register agent to registry
func (s *Socket) Register(r Registrator) {
	s.commands = r.RegisterCommands()

	syscall.Unlink(s.unixSock)
	ln, err := net.Listen("unix", s.unixSock)
	s.ready()
	if err != nil {
		panic(err)
	}

	for {
		fd, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go s.resolver(fd)
	}
}

//Command sends command to socket and gets output
func (s *Socket) Command(cmd string) {
	c, err := net.Dial("unix", s.unixSock)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	_, err = c.Write([]byte(cmd))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 512)
	for {
		n, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}

			panic(err)
		}

		fmt.Printf("%s", string(buf[0:n]))
	}

}

// resolve sent commands to socket
// and basically resolve that
func (s Socket) resolver(c net.Conn) {
	defer c.Close()
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		req := strings.Fields(strings.Trim(string(data), "\n"))
		if len(req) < 2 {
			c.Write([]byte("Invalid command"))
			return
		}

		command := req[1]
		args := req[2:]

		if command != "" {
			r, err := s.Call(c, command, args...)
			if err != nil {
				c.Write([]byte(err.Error()))
			}

			table := tablewriter.NewWriter(c)
			for _, v := range r {
				table.Append(v)
			}

			table.Render() // Send output
			return
		}
	}
}

// Call specified command
func (s Socket) Call(r io.Writer, name string, args ...string) ([][]string, error) {
	if command, exist := s.commands[name]; exist {
		return command.Closure(r, args...)
	}

	return nil, errors.New("Command does not exists, or it's not registered")
}

// New creates new socket
func New(socketPath string) *Socket {
	sock := &Socket{
		unixSock:  socketPath,
		commands:  make(map[string]Command),
		readyChan: make(chan bool, 1),
	}

	return sock
}

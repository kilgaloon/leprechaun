package api

import (
	"fmt"
	"io"
	"net"
	"strings"
	"syscall"

	"github.com/olekukonko/tablewriter"
)

// Socket defines socket on which we listen for commands
type Socket struct {
	unixSock string
	registry map[string]*Registrator
}

// This interface implements method
// that agents will use to register for command socket
type registry interface {
	RegisterCommandSocket() *Registrator
}

// Register agent to registry
func (s *Socket) Register(r registry) {
	reg := r.RegisterCommandSocket()
	s.registry[reg.Agent.GetName()] = reg

	syscall.Unlink(s.unixSock)
	ln, err := net.Listen("unix", s.unixSock)
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
			continue
		}

		registry := req[0]
		command := req[1]
		args := req[2:]

		if registry != "" || command != "" {
			r, err := s.registry[registry].Call(command, args...)
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

// BuildSocket create new socket instance
func BuildSocket(socketPath string) *Socket {
	return &Socket{
		unixSock: socketPath,
		registry: make(map[string]*Registrator),
	}
}

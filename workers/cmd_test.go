package workers

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"io"
	syslog "log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

var (
	iniFile         = "../tests/configs/config_regular.ini"
	path            = &iniFile
	cfgWrap         = config.NewConfigs()
	remoteRecipe, _ = recipe.Build("../tests/etc/leprechaun/recipes/remote.yml")
	wrks            = New(
		cfgWrap.New("test", *path),
		log.Logs{},
		context.New(),
		true,
	)
	wrksNotDebug = New(
		cfgWrap.New("test", *path),
		log.Logs{},
		context.New(),
		false,
	)
	remoteWorker, _         = wrks.CreateWorker(remoteRecipe)
	remoteWorkerNotDebug, _ = wrksNotDebug.CreateWorker(remoteRecipe)
)

func TestMain(t *testing.T) {
	startTCP(func() {
		remoteWorker.Run()
	})

	startTLSTCP(func() {
		remoteWorkerNotDebug.Run()
	})

	// cleanup
	os.Remove("remote.txt")
	os.Remove("buffer_test.txt")

}

func startTCP(f func()) {
	ln, err := net.Listen("tcp", ":11400")
	if err != nil {
		syslog.Fatal(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}

			go handleConnection(conn)
		}

	}()

	time.Sleep(5 * time.Second)

	f()

	ln.Close()

}

func startTLSTCP(f func()) {
	cert, err := tls.LoadX509KeyPair(
		"../tests/crts/certificate.pem",
		"../tests/crts/key.pem",
	)

	if err != nil {
		syslog.Fatal(err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.NoClientCert}
	config.Rand = rand.Reader

	ln, err := tls.Listen("tcp", ":11400", &config)
	if err != nil {
		syslog.Fatal(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}

			go handleConnection(conn)
		}

	}()

	time.Sleep(5 * time.Second)

	f()

	ln.Close()
}

func handleConnection(c net.Conn) {
	bSize := 1024
	// tell client that server is ready to read instructions
	_, err := c.Write([]byte("ready"))
	if err != nil {
		c.Close()
	}

	var input []byte
	var n int
	var cont string

	input = make([]byte, bSize)
	n, err = c.Read(input)
	if err != nil {
		c.Close()
	}

	cont += string(input)
	// if buffer size and read size are equal
	// that means that there is more to read from socket
	for n == bSize {
		input = make([]byte, bSize)
		n, err = c.Read(input)
		if err != nil {
			c.Close()
		}

		if n == 0 {
			break
		}

		bytes.Trim(input, "\x00")
		cont += string(input)
	}

	var b bytes.Buffer
	_, err = b.Write([]byte(cont))
	if err != nil && err != io.EOF {
		c.Close()
	}

	_, err = c.Write([]byte(">"))
	if err != nil {
		c.Close()
	}

	cmd := make([]byte, 256)

	_, err = c.Read(cmd)
	if err != nil {
		c.Close()
	}

	//s := Step(string(bytes.Trim(cmd, "\x00")))
	// if s.Validate() {
	// 	cmd, err := NewCmd(s, &b, r.Context, r.Debug, "native")
	// 	if err != nil {
	// 		_, err = c.Write([]byte("error"))
	// 	}

	// 	if r.isCmdAllowed(cmd) {
	// 		err = cmd.Run()
	// 		if err != nil {
	// 			r.Error(err.Error())
	// 			_, err = c.Write([]byte("error"))
	// 		}

	// 		_, err = c.Write(cmd.Stdout.Bytes())
	// 		if err != nil {
	// 			r.Error(err.Error())
	// 		}
	// 	} else {
	// 		r.Error("Command not allowed %s", s.Name())
	// 		_, err = c.Write([]byte("error"))
	// 	}
	// }

	c.Close()
}

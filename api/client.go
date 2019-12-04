package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/olekukonko/tablewriter"
)

// Cmd is type of command that is passed to resolver
type Cmd string

// Agent returns to which agent is command refered to
func (c Cmd) Agent() string {
	s := strings.Fields(string(c))

	if len(s) > 0 {
		return s[0]
	}

	return ""
}

// Command is what command agent needs to execute
func (c Cmd) Command() string {
	s := strings.Fields(string(c))
	if len(s) > 1 {
		return s[1]
	}

	return ""
}

// Args is which args are passed to command
func (c Cmd) Args() []string {
	s := strings.Fields(string(c))
	if len(s) > 2 {
		return s[2:]
	}

	return []string{""}
}

const (
	host = "http://localhost:11401"

	infoEndpoint        = "/{agent}/info"
	workersListEndpoint = "/{agent}/workers/list"
	workersKillEndpoint = "/{agent}/workers/kill"
)

var (
	table = tablewriter.NewWriter(os.Stdout)
	// HTTPClient config
	HTTPClient  = &http.Client{Timeout: 30 * time.Second}
	processCmds = map[string]string{
		"start": "/{agent}/start",
		"stop":  "/{agent}/stop",
		"pause": "/{agent}/pause",
	}
)

// Resolver has job to resolve which enpoint to ping and return information
func Resolver(c Cmd) {
	e := host + "/" + c.Agent() + "/"
	cmd := strings.Replace(c.Command(), ":", "/", -1)
	endpoint, _ := url.Parse(e + cmd)

	q := endpoint.Query()
	if len(c.Args()) > 0 {
		for _, arg := range c.Args() {
			q.Add("args", arg)
		}
	}

	fullEndpoint := endpoint.String()
	if len(q) > 0 {
		fullEndpoint += "?" + q.Encode()
	}

	r, err := HTTPClient.Get(fullEndpoint)
	if err != nil && err != io.EOF {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	if r.StatusCode == 404 {
		fmt.Println("No such command")
		return
	}

	defer r.Body.Close()

	resp := &TableResponse{}
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	table.SetHeader(resp.Header)

	for _, col := range resp.Columns {
		table.Append(col)
	}

	table.Render()

}

// RevealEndpoint formats endpoint to be used for interal http api
func RevealEndpoint(e string, c Cmd) string {
	return host + strings.Replace(e, "{agent}", c.Agent(), -1)
}

// IsAPIRunning checks is http api running
func IsAPIRunning() bool {
	_, err := HTTPClient.Get(host)
	if err != nil {
		return false
	}

	return true
}

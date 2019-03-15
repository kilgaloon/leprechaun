package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Cmd is type of command that is passed to resolver
type Cmd string

// Agent returns to which agent is command refered to
func (c Cmd) Agent() string {
	s := strings.Fields(string(c))
	return s[0]
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
	switch c.Command() {
	case "info":
		Info(c)
		break
	case "workers:list":
		WorkersList(c)
		break
	case "workers:kill":
		WorkersKill(c)
		break
	default:
		Process(c)
		break
	}
}

// RevealEndpoint formats endpoint to be used for interal http api
func RevealEndpoint(e string, c Cmd) string {
	return host + strings.Replace(e, "{agent}", c.Agent(), -1)
}

// Info display info for the agent
func Info(c Cmd) {
	r, err := HTTPClient.Get(RevealEndpoint(infoEndpoint, c))
	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode != 200 {
		fmt.Println("No such command")
		return
	}

	defer r.Body.Close()

	resp := &InfoResponse{}
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		log.Fatal(err)
	}

	table.SetHeader([]string{"Status", "Recipes in queue"})
	table.Append([]string{resp.Status, resp.RecipesInQueue})

	table.Render()
}

// Process command handles those endpoints that start, stop and pause
func Process(c Cmd) {
	r, err := HTTPClient.Get(RevealEndpoint(processCmds[c.Command()], c))
	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode != 200 {
		fmt.Println("No such command " + c.Command())
		return
	}

	defer r.Body.Close()

	resp := &WorkersResponse{}
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		log.Fatal(err)
	}

	table.Append([]string{resp.Message})
	table.Render()
}

// WorkersList display info for the agent
func WorkersList(c Cmd) {
	r, err := HTTPClient.Get(RevealEndpoint(workersListEndpoint, c))
	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode != 200 {
		fmt.Println("No such command")
		return
	}

	defer r.Body.Close()

	resp := &WorkersResponse{}
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Message == "" {
		table.SetHeader([]string{"Name", "Started at", "Working on"})
		for _, w := range resp.List {
			table.Append(w)
		}
	} else {
		table.Append([]string{resp.Message})
	}

	table.Render()
}

// WorkersKill display info for the agent
func WorkersKill(c Cmd) {
	r, err := HTTPClient.Get(RevealEndpoint(workersKillEndpoint, c) + "?name=" + c.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode != 200 {
		fmt.Println("No such command")
		return
	}

	defer r.Body.Close()

	resp := &WorkersResponse{}
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		log.Fatal(err)
	}

	table.Append([]string{resp.Message})

	table.Render()
}

// IsAPIRunning checks is http api running
func IsAPIRunning() bool {
	_, err := HTTPClient.Get(host)
	if err != nil {
		return false
	}

	return true
}

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

func (c Cmd) agent() string {
	s := strings.Fields(string(c))
	return s[0]
}

func (c Cmd) command() string {
	s := strings.Fields(string(c))
	return s[1]
}

func (c Cmd) args() []string {
	s := strings.Fields(string(c))
	return s[2:]
}

const (
	host = "http://localhost:11401"

	infoEndpoint        = "/{agent}/info"
	workersListEndpoint = "/{agent}/workers/list"
	workersKillEndpoint = "/{agent}/workers/kill"
)

var (
	table      = tablewriter.NewWriter(os.Stdout)
	httpClient = &http.Client{Timeout: 30 * time.Second}
)

// Resolver has job to resolve which enpoint to ping and return information
func Resolver(c Cmd) {
	switch c.command() {
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
		fmt.Println("No such command")
	}
}

func revealEndpoint(e string, c Cmd) string {
	return host + strings.Replace(e, "{agent}", c.agent(), -1)
}

// Info display info for the agent
func Info(c Cmd) {
	r, err := httpClient.Get(revealEndpoint(infoEndpoint, c))
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

	table.SetHeader([]string{"PID", "Config file", "Recipes in queue", "Memory allocated"})
	table.Append([]string{resp.PID, resp.ConfigFile, resp.RecipesInQueue, resp.MemoryAllocated})

	table.Render()
}

// WorkersList display info for the agent
func WorkersList(c Cmd) {
	r, err := httpClient.Get(revealEndpoint(workersListEndpoint, c))
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
	r, err := httpClient.Get(revealEndpoint(workersKillEndpoint, c) + "?name=" + c.args()[0])
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
	_, err := httpClient.Get(host)
	if err != nil {
		return false
	}

	return true
}

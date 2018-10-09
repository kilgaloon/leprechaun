package workers

import (
	"testing"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
)

var (
	workers = New(
		1,
		"../tests/var/log/leprechaun/workers.output",
		log.Logs{},
		context.New(),
	)
)

func TestCreateWorker(t *testing.T) {
	_, err := workers.CreateWorker("test")
	if err != nil {
		t.Fail()
	}

	_, err = workers.CreateWorker("test")
	if err == nil {
		t.Fail()
	}

	// test that size can't be more then 1
	_, err = workers.CreateWorker("test2")
	if err == nil {
		t.Fail()
	}
}

func TestGetWorkerByName(t *testing.T) {
	_, err := workers.GetByName("test")
	if err != nil {
		t.Fail()
	}

	_, err = workers.GetByName("test2")
	if err == nil {
		t.Fail()
	}
}

func TestGetAll(t *testing.T) {
	w := workers.GetAll()
	if workers.Size() != len(w) {
		t.Fail()
	}
}

func TestWorkerIsDone(t *testing.T) {
	workers.DoneChan <- "test"
}

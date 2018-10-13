package workers

import (
	"testing"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
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
	r, err := recipe.Build("../tests/etc/leprechaun/recipes/schedule.yml")

	_, err = workers.CreateWorker(&r)
	if err != nil {
		t.Fail()
	}

	_, err = workers.CreateWorker(&r)
	if err == nil {
		t.Fail()
	}

	// test that size can't be more then 1
	_, err = workers.CreateWorker(&r)
	if err == nil {
		t.Fail()
	}
}

func TestGetWorkerByName(t *testing.T) {
	r, err := recipe.Build("../tests/etc/leprechaun/recipes/schedule.yml")
	workers.CreateWorker(&r)

	_, err = workers.GetWorkerByName("schedule")
	if err != nil {
		t.Fail()
	}

	_, err = workers.GetWorkerByName("test2")
	if err == nil {
		t.Fail()
	}
}

func TestGetAll(t *testing.T) {
	w := workers.GetAllWorkers()
	if workers.NumOfWorkers() != len(w) {
		t.Fail()
	}
}

func TestWorkerIsDone(t *testing.T) {
	workers.DoneChan <- "schedule"
}

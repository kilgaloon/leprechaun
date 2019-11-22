package workers

import (
	"testing"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

var (
	workers = New(
		ConfigWithSettings,
		log.Logs{},
		context.New(),
	)

	workers5 = New(
		ConfigWithSettings,
		log.Logs{},
		context.New(),
	)

	workers4 = New(
		ConfigWithSettings,
		log.Logs{},
		context.New(),
	)

	workers3 = New(
		ConfigWithQueueSettings,
		log.Logs{},
		context.New(),
	)
)

func TestCreateWorker(t *testing.T) {
	r, err := recipe.Build("../tests/etc/leprechaun/recipes/schedule.yml")
	if err != nil {
		t.Fail()
	}

	_, err = workers.CreateWorker(r)
	if err != nil {
		t.Fail()
	}

	_, err = workers.CreateWorker(r)
	if err == nil {
		t.Fail()
	}
}

func TestKillWorker(t *testing.T) {
	r, err := recipe.Build("../tests/etc/leprechaun/recipes/sleep.yml")
	if err != nil {
		t.Fatal(err)
	}

	w, err := workers.CreateWorker(r)
	if err != nil {
		t.Fatal(err)
	}

	go w.Run()
	w.Kill()
}

func TestCreateWorkerQueue(t *testing.T) {
	r, _ := recipe.Build("../tests/etc/leprechaun/recipes/schedule.yml")
	r2, _ := recipe.Build("../tests/etc/leprechaun/recipes/schedule.2.yml")
	r3, _ := recipe.Build("../tests/etc/leprechaun/recipes/schedule.3.yml")

	workers3.CreateWorker(r)
	workers3.CreateWorker(r2)

	_, err := workers3.CreateWorker(r3)
	if err == nil {
		t.Fail()
	}

	workers3.Queue.pop()
	if workers3.Queue.len() > 0 {
		t.Fail()
	}
}

func TestGetWorkerByName(t *testing.T) {
	r, err := recipe.Build("../tests/etc/leprechaun/recipes/schedule.yml")
	if err != nil {
		t.Fail()
	}

	workers.CreateWorker(r)
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

	workers.Lock()
	l := len(w)
	workers.Unlock()

	if workers.NumOfWorkers() != l {
		t.Fail()
	}
}

func TestWorkerIsDone(t *testing.T) {
	workers.DoneChan <- "schedule"
}

func TestWorkerError(t *testing.T) {
	r, _ := recipe.Build("../tests/etc/leprechaun/recipes/invalid.yml")
	w, _ := workers.CreateWorker(r)

	go w.Run()
}

func TestWorkerAsyncTask(t *testing.T) {
	r, _ := recipe.Build("../tests/etc/leprechaun/recipes/schedule.3.yml")
	w, _ := workers4.CreateWorker(r)

	w.Run()
}

func TestWorkerPipeTask(t *testing.T) {
	r, _ := recipe.Build("../tests/etc/leprechaun/recipes/schedule.4.yml")
	w, _ := workers5.CreateWorker(r)

	w.Run()
}

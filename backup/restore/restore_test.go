package restore_test

import (
	"testing"

	"github.com/stefanprodan/mgob/backup/restore"
	"github.com/stefanprodan/mgob/config"
)

func assertError(t *testing.T, err error) {
	t.Log(err)
	if err == nil {
		t.Error(err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Log(err)
	if err != nil {
		t.Error(err)
	}
}

func TestMongoRestoreReturnErrorOnInvalidHost(t *testing.T) {
	target := config.Target{
		Host: "invalid",
		Port: 27017,
	}
	sched := config.Scheduler{Timeout: 60}
	plan := config.Plan{
		Target:    target,
		Scheduler: sched,
	}
	err := restore.Restore(plan, "invalid")
	assertError(t, err)
}

func TestMongoRestoreReturnErrorOnInvalidArchive(t *testing.T) {
	target := config.Target{
		Host: "localhost",
		Port: 27017,
	}
	sched := config.Scheduler{Timeout: 60}
	plan := config.Plan{
		Target:    target,
		Scheduler: sched,
	}
	err := restore.Restore(plan, "invalid")
	assertError(t, err)
}

func setUp(localhost string, port int) config.Plan {
	target := config.Target{
		Host: "localhost",
		Port: 27017,
	}
	sched := config.Scheduler{Timeout: 60}
	plan := config.Plan{
		Target:    target,
		Scheduler: sched,
	}
	return plan
}

func TestMongoRestoreWithSucess(t *testing.T) {
	plan := setUp("localhost", 27017)
	err := restore.Restore(plan, "/tmp/dump_test.gz")
	assertNoError(t, err)
}

// +build unit

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

func TestMongoRestoreReturnErrorOnInvalidArchive(t *testing.T) {
	sched := config.Scheduler{Timeout: 60}
	plan := config.Plan{
		Target:    config.Target{},
		Scheduler: sched,
	}
	_, err := restore.Restore(plan, "invalid")
	assertError(t, err)
}

func setUp(collCount int, collsLength int) config.Plan {
	collections := []config.Collection{
		{
			Name:  "parameters",
			Count: collCount,
		},
	}
	restore := config.Restore{
		Database:          "garden",
		Collections:       collections,
		CollectionsLength: collsLength,
	}
	sched := config.Scheduler{Timeout: 60}
	plan := config.Plan{
		Scheduler: sched,
		Restore:   restore,
	}
	return plan
}

func TestMongoRestoreWithSuccess(t *testing.T) {
	collCount := 5
	collsLength := 1
	plan := setUp(collCount, collsLength)
	_, err := restore.Restore(plan, "/tmp/dump_test.gz")
	assertNoError(t, err)
}

func TestShouldGetErrorOnInvalidCount(t *testing.T) {
	collCount := 10
	collsLength := 1
	plan := setUp(collCount, collsLength)
	_, err := restore.Restore(plan, "/tmp/dump_test.gz")
	assertError(t, err)
}

func TestShouldGetErrorOnInvalidCollectionLength(t *testing.T) {
	collCount := 5
	collsLength := 2
	plan := setUp(collCount, collsLength)
	_, err := restore.Restore(plan, "/tmp/dump_test.gz")
	assertError(t, err)
}

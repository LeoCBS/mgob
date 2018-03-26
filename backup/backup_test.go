// +build integration

package backup_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stefanprodan/mgob/backup"
	"github.com/stefanprodan/mgob/config"
	"gopkg.in/mgo.v2"
)

func assertError(t *testing.T, err error) {
	t.Log(err)
	if err == nil {
		t.Error(err)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func setUp(host string, port int, restore config.Restore) config.Plan {
	target := config.Target{
		Host: host,
		Port: port,
	}
	plan := config.Plan{
		Target:  target,
		Restore: restore,
	}
	return plan
}

func TestShouldGetErrorOnInvalidHost(t *testing.T) {
	restore := config.Restore{
		Database:   "whatever",
		Collection: "parameters",
		Count:      5,
	}
	plan := setUp("invalid.com", 0, restore)
	tmpPath := "/tmp"
	storagePath := "/storage"
	res, err := backup.Run(plan, tmpPath, storagePath)
	t.Logf("backup result %v", res)
	assertError(t, err)
}
func TestShouldGetErrorOnInvalidDatabase(t *testing.T) {
	port, err := strconv.Atoi(os.Getenv("MONGODB_PORT_27017_TCP_PORT"))
	assertNoError(t, err)
	host := os.Getenv("MONGODB_PORT_27017_TCP_ADDR")
	restore := config.Restore{
		Database:   "wrong",
		Collection: "parameters",
		Count:      5,
	}
	plan := setUp(host, port, restore)
	tmpPath := "/tmp"
	storagePath := "/storage"
	res, err := backup.Run(plan, tmpPath, storagePath)
	t.Logf("backup result %v", res)
	assertError(t, err)
}

func TestShouldRunBackupCorrectly(t *testing.T) {
	port, err := strconv.Atoi(os.Getenv("MONGODB_PORT_27017_TCP_PORT"))
	assertNoError(t, err)
	host := os.Getenv("MONGODB_PORT_27017_TCP_ADDR")
	restore := config.Restore{
		Database:   "garden",
		Collection: "parameters",
		Count:      15,
	}
	plan := setUp(host, port, restore)
	s, err := getMongoSession(host, port)
	assertNoError(t, err)
	docsLength := 20
	insertMongoData(t, s, restore.Database, restore.Collection, docsLength)
	defer tearDown(t, s, restore.Database, restore.Collection)
	tmpPath := "/tmp"
	storagePath := "/storage"
	res, err := backup.Run(plan, tmpPath, storagePath)
	t.Logf("backup result %v", res)
	assertNoError(t, err)
}

func tearDown(t *testing.T, s *mgo.Session, database, collection string) {
	names, err := s.DatabaseNames()
	assertNoError(t, err)
	for _, name := range names {
		err = s.DB(name).DropDatabase()
		assertNoError(t, err)
	}
}

func getMongoSession(host string, port int) (*mgo.Session, error) {
	mongoURL := fmt.Sprintf("mongodb://%v:%d", host, port)
	return mgo.Dial(mongoURL)
}

func insertMongoData(t *testing.T, s *mgo.Session, database, collection string, length int) {
	c := s.DB(database).C(collection)
	type parameter struct{ count int }
	sum := 0
	for i := 0; i < length; i++ {
		err := c.Insert(&parameter{i})
		assertNoError(t, err)
		sum += i
	}
}

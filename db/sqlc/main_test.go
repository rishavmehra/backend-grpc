package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:hide1337@localhost:5432/simple_bank?sslmode=disable"
)

var testQuries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Not able to make connection with Database.", err)
	}

	testQuries = New(testDB)
	os.Exit(m.Run())
}

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

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Not able to make connection with Database.", err)
	}

	testQuries = New(conn)
	os.Exit(m.Run())
}

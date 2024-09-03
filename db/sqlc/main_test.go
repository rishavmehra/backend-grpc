package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rishavmehra/backend-grpc/util"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:hide1337@localhost:5432/simple_bank_new?sslmode=disable"
// )

var testQuries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot Load configuration: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Not able to make connection with Database.", err)
	}

	testQuries = New(testDB)
	os.Exit(m.Run())
}

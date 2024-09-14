package main

import (
	"database/sql"
	"log"

	"github.com/rishavmehra/backend-grpc/api"
	db "github.com/rishavmehra/backend-grpc/db/sqlc"
	"github.com/rishavmehra/backend-grpc/util"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://root:hide1337@localhost:5432/grpc_staging?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot Load configuration: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Not able to make connection with Database.", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot Create Server", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can't start the server")
	}
}

package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rishavmehra/backend-grpc/api"
	db "github.com/rishavmehra/backend-grpc/db/sqlc"
	"github.com/rishavmehra/backend-grpc/gapi"
	pb "github.com/rishavmehra/backend-grpc/pb"
	"github.com/rishavmehra/backend-grpc/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

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
	go gatewayServer(config, store)
	grpcServer(config, store)

}

func grpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot Create Server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot start the server", err)
	}

	log.Printf("Starting GRPC server on %s", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Cannot start the server", err)
	}
}

func gatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot Create Server", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot start the server", err)
	}

	log.Printf("Starting HTTP gateway server on %s", listener.Addr().String())
	if err := http.Serve(listener, mux); err != nil {
		log.Fatal("Cannot start the HTTP Gateway server", err)
	}
}

func ginServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot Create Server", err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Can't start the server")
	}
}

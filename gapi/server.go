package gapi

import (
	"fmt"

	db "github.com/rishavmehra/backend-grpc/db/sqlc"
	pb "github.com/rishavmehra/backend-grpc/pb"
	"github.com/rishavmehra/backend-grpc/token"
	"github.com/rishavmehra/backend-grpc/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot Create Token Maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

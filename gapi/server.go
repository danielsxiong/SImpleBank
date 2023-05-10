package gapi

import (
	db "danielsxiong/simplebank/db/sqlc"
	"danielsxiong/simplebank/pb"
	"danielsxiong/simplebank/token"
	"danielsxiong/simplebank/util"
	"fmt"
)

// Server serve http requests for banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new grpc server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPASETOMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

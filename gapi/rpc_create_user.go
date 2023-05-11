package gapi

import (
	"context"
	db "danielsxiong/simplebank/db/sqlc"
	"danielsxiong/simplebank/pb"
	"danielsxiong/simplebank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err.Error())
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err.Error())
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err.Error())
	}

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return resp, nil
}

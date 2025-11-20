package user

import (
	"context"
	"errors"
	"user/internal/storage"

	userpb "github.com/steephseqq/maximprotos-user/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceGRPC interface {
	CreateUser(
		ctx context.Context,
		username, name, email, uuid, avatarURL string,
	) (err error)
	DeleteUser(
		ctx context.Context,
		uuid string,
	) (err error)
}

type serverAPI struct {
	userpb.UnimplementedUserServer
	user UserServiceGRPC
}

func Register(gRPC *grpc.Server, userService UserServiceGRPC) {
	userpb.RegisterUserServer(gRPC, &serverAPI{user: userService})
}

func (s *serverAPI) CreateUser(
	ctx context.Context,
	req *userpb.CreateUserRequest,
) (*userpb.CreateUserResponse, error) {
	var (
		uuid      = req.GetUuid()
		email     = req.GetEmail()
		username  = req.GetUsername()
		name      = req.GetName()
		avatarURL = "https://i.pinimg.com/736x/95/73/95/957395e49e0e1ea58efc58c7159778e8.jpg"
	)
	if email == "" || username == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	if err := s.user.CreateUser(ctx, username, name, email, uuid, avatarURL); err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userpb.CreateUserResponse{
		UUID: uuid,
	}, nil
}

func (s *serverAPI) DeleteUser(
	ctx context.Context,
	req *userpb.DeleteUserRequest,
) (*userpb.DeleteUserResponse, error) {
	if req.UUID == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}

	if err := s.user.DeleteUser(ctx, req.UUID); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userpb.DeleteUserResponse{
		Success: true,
	}, nil
}

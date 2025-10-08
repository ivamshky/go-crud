package handler

import (
	"context"

	pb "github.com/ivamshky/go-crud/gen/grpc/user"
	"github.com/ivamshky/go-crud/model"
	"github.com/ivamshky/go-crud/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func (server *UserServer) GetUser(context context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	userDetails, err := server.userService.GetUserDetails(context, req.Id)
	if err != nil {
		return nil, err
	}
	return convertFromUserModelToPb(userDetails), nil
}

func (server *UserServer) CreateUser(context context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	userDetails, err := server.userService.CreateUser(context, model.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	})
	if err != nil {
		return nil, err
	}

	return convertFromUserModelToPb(userDetails), nil
}

func (server *UserServer) UpdateUser(context context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func (server *UserServer) DeleteUser(context context.Context, req *pb.DeleteUserRequest) (*pb.EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func convertFromUserModelToPb(model model.User) *pb.User {
	return &pb.User{
		Id:    model.Id,
		Name:  model.Name,
		Email: model.Email,
		Age:   model.Age,
	}
}

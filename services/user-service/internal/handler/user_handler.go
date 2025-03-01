package handler

import (
	"context"
	model "github.com/ppthinh/ChatApp/services/user-service/internal/models"
	pb "github.com/ppthinh/ChatApp/services/user-service/internal/proto/genproto"
	"github.com/ppthinh/ChatApp/services/user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type userGrpcServer struct {
	svc service.UserService
	pb.UnimplementedUserServiceServer
}

func RunUserGRPCServer(svc service.UserService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Can't listen gRPC: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userGrpcServer{svc: svc})
	reflection.Register(s)
	log.Printf("gRPC Server start at port: %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Can't start gRPC server: %v", err)
	}
}

func (u *userGrpcServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.User, error) {
	user, err := u.svc.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (u *userGrpcServer) GetUserByPhoneNumber(ctx context.Context, req *pb.GetUserByPhoneNumberRequest) (*pb.User, error) {
	user, err := u.svc.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (u *userGrpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := u.svc.CreateUser(req)
	if err != nil {
		return &pb.CreateUserResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.CreateUserResponse{
		UserId:  user.ID.String(),
		Success: true,
	}, nil

}

func (u *userGrpcServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := u.svc.UpdateUser(req)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func (u *userGrpcServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	//userId, _ := uuid.Parse(req.Id)
	if err := u.svc.DeleteUser(req.UserId); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func mapUserToResponse(user *model.User) *pb.User {
	return &pb.User{
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
}

func (u *userGrpcServer) GetUserForAuth(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserForAuth, error) {
	user, err := u.svc.GetUserForAuth(req.Email)
	if err != nil {
		return nil, err
	}

	userID, err := user.ID.MarshalText()
	if err != nil {
		return nil, err
	}

	return &pb.UserForAuth{
		UserId:   string(userID),
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

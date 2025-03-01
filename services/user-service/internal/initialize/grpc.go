package initialize

import (
	"context"
	"github.com/google/uuid"
	model "github.com/ppthinh/ChatApp/services/user-service/internal/models"
	pb "github.com/ppthinh/ChatApp/services/user-service/internal/proto/genproto"
	"github.com/ppthinh/ChatApp/services/user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type grpcServer struct {
	svc service.UserService
	pb.UnimplementedUserServiceServer
}

func runGRPCServer(svc service.UserService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Can't listen gRPC: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &grpcServer{svc: svc})
	reflection.Register(s)
	log.Printf("gRPC Server start at port: %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Can't start gRPC server: %v", err)
	}
}

func (s *grpcServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.User, error) {
	user, err := s.svc.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (s *grpcServer) GetUserByPhoneNumber(ctx context.Context, req *pb.GetUserByPhoneNumberRequest) (*pb.User, error) {
	user, err := s.svc.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (s *grpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user := &model.User{
		ID:          uuid.New(),
		Name:        req.Name,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}
	if err := s.svc.CreateUser(user); err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (s *grpcServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user := &model.User{
		Name:        req.Name,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}
	if err := s.svc.UpdateUser(user); err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (s *grpcServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	id, _ := uuid.Parse(req.Id)
	if err := s.svc.DeleteUser(id); err != nil {
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

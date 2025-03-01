package handler

import (
	"context"
	pb "github.com/ppthinh/ChatApp/services/auth-service/internal/proto/genproto"
	"github.com/ppthinh/ChatApp/services/auth-service/internal/service"
)

type AuthGrpcServer struct {
	pb.UnimplementedAuthServiceServer
	svc *service.AuthService
}

func NewAuthGrpcServer(svc *service.AuthService) *AuthGrpcServer {
	return &AuthGrpcServer{svc: svc}
}

func (a *AuthGrpcServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return a.svc.RegisterUser(req)
}

func (a *AuthGrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return a.svc.LoginUser(ctx, req)
}

//func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
//	userID, err := s.svc.ValidateToken(req.Token)
//	if err != nil {
//		return &pb.ValidateTokenResponse{
//			Valid:        false,
//			ErrorMessage: err.Error(),
//		}, nil
//	}
//	return &pb.ValidateTokenResponse{
//		UserId: userID,
//		Valid:  true,
//	}, nil
//}

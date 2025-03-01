package initialize

import (
	"github.com/ppthinh/ChatApp/services/auth-service/internal/handler"
	pb "github.com/ppthinh/ChatApp/services/auth-service/internal/proto/genproto"
	userpb "github.com/ppthinh/ChatApp/services/auth-service/internal/proto/genproto"
	"github.com/ppthinh/ChatApp/services/auth-service/internal/service"
	"github.com/ppthinh/ChatApp/services/auth-service/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

func Run() {
	cfg := config.LoadConfig()

	conn, err := grpc.NewClient(cfg.UserGRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer conn.Close()

	userClient := userpb.NewUserServiceClient(conn)
	authSvc := service.NewAuthService(userClient, cfg.JWTSecretKey)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, handler.NewAuthGrpcServer(authSvc))

	log.Printf(cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

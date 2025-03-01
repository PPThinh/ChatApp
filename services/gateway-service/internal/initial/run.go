package initial

import (
	"github.com/ppthinh/ChatApp/services/gateway-service/internal/handler"
	authpb "github.com/ppthinh/ChatApp/services/gateway-service/internal/proto/genproto"
	userpb "github.com/ppthinh/ChatApp/services/gateway-service/internal/proto/genproto"
	"github.com/ppthinh/ChatApp/services/gateway-service/internal/router"
	"github.com/ppthinh/ChatApp/services/gateway-service/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func newGRPCClient(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Run() {
	cfg := config.LoadConfig()

	userConn, err := newGRPCClient(cfg.GrpcUserAddress)
	if err != nil {
		log.Fatalf("failed to connect to user service at %s: %v", cfg.GrpcUserAddress, err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	authConn, err := newGRPCClient(cfg.GrpcAuthAddress)
	if err != nil {
		log.Fatalf("failed to connect to auth service at %s: %v", cfg.GrpcAuthAddress, err)
	}
	defer authConn.Close()
	authClient := authpb.NewAuthServiceClient(authConn)

	handle := handler.NewGatewayHandler(userClient, authClient)
	log.Printf("API Gateway starting on %s", ":"+cfg.Port)
	if err := router.StartGateway(":"+cfg.Port, handle); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	model "github.com/ppthinh/ChatApp/services/user-service/internal/models"
	pb "github.com/ppthinh/ChatApp/services/user-service/internal/proto/genproto"
	"github.com/ppthinh/ChatApp/services/user-service/internal/repository"
	"github.com/ppthinh/ChatApp/services/user-service/internal/service"
	"github.com/ppthinh/ChatApp/services/user-service/pkg/config"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql" // Thay sqlite bằng mysql
	"gorm.io/gorm"
	"log"
	"net"
)

type grpcServer struct {
	svc service.UserService
	pb.UnimplementedUserServiceServer
}

func Run() {
	// Load cấu hình
	cfg := config.LoadConfig()

	// Khởi tạo database MySQL
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Không thể kết nối database: %v", err)
	}
	db.AutoMigrate(&model.User{}, &model.FriendShip{})

	// Khởi tạo repository và service
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)

	// Khởi chạy gRPC server trong goroutine
	go runGRPCServer(svc, cfg.GRPCPort)

	// Khởi chạy HTTP server
	runHTTPServer(svc, cfg.Port)
}

func runGRPCServer(svc service.UserService, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Không thể lắng nghe gRPC: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &grpcServer{svc: svc})

	log.Printf("gRPC Server chạy trên cổng %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Không thể khởi động gRPC server: %v", err)
	}
}

func runHTTPServer(svc service.UserService, port string) {
	r := gin.Default()

	r.GET("/users/email/:email", func(c *gin.Context) {
		email := c.Param("email")
		user, err := svc.GetUserByEmail(email)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, mapUserToResponse(user))
	})

	r.GET("/users/phone/:phone", func(c *gin.Context) {
		phone := c.Param("phone")
		user, err := svc.GetUserByPhoneNumber(phone)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, mapUserToResponse(user))
	})

	r.POST("/users", func(c *gin.Context) {
		var req pb.CreateUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}
		user := &repository.User{
			ID:          uuid.New(),
			Name:        req.Name,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
		}
		if err := svc.CreateUser(user); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, mapUserToResponse(user))
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		id, _ := uuid.Parse(c.Param("id"))
		var req pb.UpdateUserRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}
		user := &repository.User{
			ID:          id,
			Name:        req.Name,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
		}
		if err := svc.UpdateUser(user); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, mapUserToResponse(user))
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		id, _ := uuid.Parse(c.Param("id"))
		if err := svc.DeleteUser(id); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Xóa user thành công"})
	})

	log.Printf("HTTP Server chạy trên cổng %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Không thể khởi động HTTP server: %v", err)
	}
}

// gRPC method implementations
func (s *grpcServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	user, err := s.svc.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (s *grpcServer) GetUserByPhoneNumber(ctx context.Context, req *pb.GetUserByPhoneNumberRequest) (*pb.UserResponse, error) {
	user, err := s.svc.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (s *grpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user := &repository.User{
		ID:          uuid.New(),
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}
	if err := s.svc.CreateUser(user); err != nil {
		return nil, err
	}
	return mapUserToResponse(user), nil
}

func (s *grpcServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id, _ := uuid.Parse(req.Id)
	user := &repository.User{
		ID:          id,
		Name:        req.Name,
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

func mapUserToResponse(user *repository.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
}

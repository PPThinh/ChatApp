package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ppthinh/ChatApp/services/auth-service/internal/util"
	"strings"
	"time"

	pb "github.com/ppthinh/ChatApp/services/auth-service/internal/proto/genproto"
	userpb "github.com/ppthinh/ChatApp/services/auth-service/internal/proto/genproto"
)

type registerRequest struct {
	Email       string `validate:"required,email"`
	PhoneNumber string `validate:"required,vietnam_phone"`
	Name        string `validate:"required"`
	Password    string `validate:"required"`
}

type loginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type AuthService struct {
	userClient userpb.UserServiceClient
	validator  *validator.Validate
	jwtSecret  string
}

func NewAuthService(userClient userpb.UserServiceClient, jwtSecret string) *AuthService {
	v := validator.New()
	v.RegisterValidation("vietnam_phone", validateVietnamPhone)

	return &AuthService{
		userClient: userClient,
		validator:  v,
		jwtSecret:  jwtSecret,
	}
}

func (as *AuthService) RegisterUser(req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Kiem tra hop le cua du lieu
	reqData := registerRequest{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    req.Password,
	}

	if err := as.validator.Struct(reqData); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.RegisterResponse{
			Success:      false,
			ErrorMessage: as.formatValidError(validationErrors[0]),
		}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// request create userr
	userReq := &pb.CreateUserRequest{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    req.Password,
	}
	// call user service
	userRes, err := as.userClient.CreateUser(ctx, userReq)
	if err != nil {
		return &pb.RegisterResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	if !userRes.Success {
		return &pb.RegisterResponse{
			Success:      false,
			ErrorMessage: userRes.ErrorMessage,
		}, nil
	}

	userResID, err := uuid.Parse(userRes.UserId)
	if err != nil {
		return &pb.RegisterResponse{
			Success:      false,
			ErrorMessage: "invalid user ID format",
		}, nil
	}
	// tao token
	tokenString, err := util.GenerateToken(userResID, as.jwtSecret)
	if err != nil {
		return &pb.RegisterResponse{
			Success:      false,
			ErrorMessage: "failed to generate token: " + err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		UserId:       userRes.UserId,
		Token:        tokenString,
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func (as *AuthService) LoginUser(ctx context.Context, userReq *userpb.LoginRequest) (*pb.LoginResponse, error) {
	userReqData := loginRequest{
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	if err := as.validator.Struct(userReqData); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.LoginResponse{
			Success:      false,
			ErrorMessage: as.formatValidError(validationErrors[0]),
		}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// call user service
	userRes, err := as.userClient.GetUserForAuth(ctx, &userpb.GetUserByEmailRequest{Email: userReq.Email})
	if err != nil {
		return &pb.LoginResponse{
			Success:      false,
			ErrorMessage: "invalid email or password",
		}, nil
	}

	if err := util.ComparePassword(userRes.Password, userReq.Password); err != nil {
		return &pb.LoginResponse{
			Success:      false,
			ErrorMessage: "invalid email or password",
		}, nil
	}

	userId, err := uuid.Parse(userRes.Id)
	if err != nil {
		return &pb.LoginResponse{
			Success:      false,
			ErrorMessage: "invalid user ID format",
		}, nil
	}

	tokenString, err := util.GenerateToken(userId, as.jwtSecret)
	if err != nil {
		return &pb.LoginResponse{
			Success:      false,
			ErrorMessage: "failed to generate token: " + err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Token:        tokenString,
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func validateVietnamPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	phone = strings.TrimSpace(phone)

	if strings.HasPrefix(phone, "+84") {
		if len(phone) != 12 {
			return false
		}
		phone = "0" + phone[3:]
	}

	if len(phone) != 10 || !strings.HasPrefix(phone, "0") {
		return false
	}

	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (as *AuthService) formatValidError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "vietnam_phone":
		return "Phone number must be start with 0 or +84"
	default:
		return "Invalid " + err.Field()
	}
}

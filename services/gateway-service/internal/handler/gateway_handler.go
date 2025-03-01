package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	authpb "github.com/ppthinh/ChatApp/services/gateway-service/internal/proto/genproto"
	userpb "github.com/ppthinh/ChatApp/services/gateway-service/internal/proto/genproto"
	"net/http"
)

type GatewayHandler struct {
	userClient userpb.UserServiceClient
	authClient authpb.AuthServiceClient
}

func NewGatewayHandler(userClient userpb.UserServiceClient, authClient userpb.AuthServiceClient) *GatewayHandler {
	return &GatewayHandler{
		userClient: userClient,
		authClient: authClient,
	}
}

func (h *GatewayHandler) Register(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.authClient.Register(context.Background(), &authpb.RegisterRequest{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil || !resp.Success {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id": resp.UserId,
		"token":   resp.Token,
	})
}

func (h *GatewayHandler) Login(ctx *gin.Context) {

}

func (h *GatewayHandler) GetUserByEmail(ctx *gin.Context) {

}

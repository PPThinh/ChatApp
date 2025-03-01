package model

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

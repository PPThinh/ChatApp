package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ppthinh/ChatApp/services/auth-service/internal/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func GenerateToken(userId uuid.UUID, secretKey string) (string, error) {
	claims := &model.Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ComparePassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return errors.New("invalid password")
	}
	return nil
}

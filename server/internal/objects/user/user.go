package user

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex"`
	PasswordHash string
}

type UserCreateDto struct {
	Username string
	Password string
}

type UserLoginDto struct {
	Username string
	Password string
}

type JWTClaim struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type WhoAmIResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

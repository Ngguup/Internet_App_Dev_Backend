package ds

import (
	"github.com/golang-jwt/jwt"
	"lab1/internal/app/role"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
	Role     role.Role
}
package handler

import (
	"context"
	"errors"
	"lab1/internal/app/ds"
	"lab1/internal/app/role"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

const jwtPrefix = "Bearer "

func (h *Handler) ParseAndValidateJWT(gCtx *gin.Context) (*ds.JWTClaims, error) {
	authHeader := gCtx.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, jwtPrefix) {
		return nil, errors.New("invalid authorization header")
	}

	jwtStr := authHeader[len(jwtPrefix):]

	// Проверка blacklist
	err := h.Redis.CheckJWTInBlacklist(context.Background(), jwtStr)
	if err == nil {
		return nil, errors.New("token is blacklisted")
	}
	if !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// Парсим токен
	token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(*ds.JWTClaims), nil
}

func (h *Handler) WithAuthCheck(assignedRoles ...role.Role) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		claims, err := h.ParseAndValidateJWT(gCtx)
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Проверка роли
		for _, r := range assignedRoles {
			if claims.Role == r {
				gCtx.Set("user_id", claims.UserID)
				gCtx.Set("role", int(claims.Role))
				gCtx.Next()
				return
			}
		}

		log.Printf("role %v is not allowed (allowed: %v)", claims.Role, assignedRoles)
		gCtx.AbortWithStatus(http.StatusForbidden)
	}
}


func GetUserID(ctx *gin.Context) uint {
	val, exists := ctx.Get("user_id")
	if !exists {
		log.Println("user_id not found in context")
		return 0
	}

	id, ok := val.(uint)
	if !ok {
		log.Println("invalid user_id type in context")
		return 0
	}

	return id
}



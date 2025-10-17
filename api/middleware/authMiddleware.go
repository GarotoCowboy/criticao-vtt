package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			handler.SendError(ctx, http.StatusUnauthorized, "authorization header not found")
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			handler.SendError(ctx, http.StatusUnauthorized, "Invalid authorization header format. Usage: Bearer <token>")
			ctx.Abort()
			return
		}
		tokenString := tokenParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return config.JWT_SECRET, nil
		})
		if err != nil {
			handler.SendError(ctx, http.StatusUnauthorized, "Invalid token or expired")
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userIDFLoat, ok := claims["user_id"].(float64); ok {
				ctx.Set("user_id", uint(userIDFLoat))
			} else {
				handler.SendError(ctx, http.StatusUnauthorized, "Invalid 'user_id' claim in token")
				ctx.Abort()
				return
			}
		} else {
			handler.SendError(ctx, http.StatusUnauthorized, "Invalid token")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

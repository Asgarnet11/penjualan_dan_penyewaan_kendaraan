package middleware

import (
	"errors"
	"net/http"
	"strings"
	"sultra-otomotif-api/internal/helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			helper.ErrorResponse(c, "Authorization header is required", http.StatusUnauthorized, errors.New("missing bearer token"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			helper.ErrorResponse(c, "Invalid token", http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helper.ErrorResponse(c, "Failed to parse claims", http.StatusUnauthorized, errors.New("invalid claims format"))
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			helper.ErrorResponse(c, "Invalid user ID in token", http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		// Set data user ke context agar bisa diakses oleh handler selanjutnya
		c.Set("currentUserID", userID)
		c.Set("currentUserRole", claims["role"].(string))
		c.Next()
	}
}

// Middleware untuk mengecek role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("currentUserRole")
		if !exists || userRole.(string) != requiredRole {
			helper.ErrorResponse(c, "You are not authorized to perform this action", http.StatusForbidden, errors.New("insufficient privileges"))
			c.Abort()
			return
		}
		c.Next()
	}
}

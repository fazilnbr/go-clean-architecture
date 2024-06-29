package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const USER_ID = "user_id"

type MiddlewareIf interface {
	CreateToken() gin.HandlerFunc
	ValidateToken() gin.HandlerFunc
}

type Middleware struct {
	secretKey string
}

func NewMiddleware(secretKey string) MiddlewareIf {
	return &Middleware{secretKey: secretKey}
}

func (m *Middleware) CreateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		id, err := c.Cookie(USER_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create token failed with " + err.Error()})
			c.Abort()
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			USER_ID: id,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		})

		tokenString, err := token.SignedString([]byte(m.secretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create token failed with " + err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user login successfully", "token": tokenString})

	}
}

func (m *Middleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userID, ok := claims[USER_ID].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user ID in the header
		c.Request.Header.Set(USER_ID, fmt.Sprintf("%d", uint(userID)))

		c.Next()
	}
}

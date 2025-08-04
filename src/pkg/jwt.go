package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

var JwtSecretKey = []byte(viper.GetString("SERVICE_SECRET"))

type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, role string, expiryMinutes int) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiryMinutes))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func getUserID(c *gin.Context) string {
	if val, exists := c.Get("user_id"); exists {
		if uid, ok := val.(string); ok {
			return uid
		}
	}
	return ""
}

func getUserRole(c *gin.Context) string {
	if val, exists := c.Get("role"); exists {
		if role, ok := val.(string); ok {
			return role
		}
	}
	return ""
}

func ExtractToken(c *gin.Context) *CustomClaims {
	return &CustomClaims{
		UserID: getUserID(c),
		Role:   getUserRole(c),
	}
}

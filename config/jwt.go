package config

import (

	"time"

	"github.com/golang-jwt/jwt/v5"
	
)

var jwtSecret = []byte("secretkey")

// GenerateJWT creates a JWT token for authentication
func GenerateJWT(userID string, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 2).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

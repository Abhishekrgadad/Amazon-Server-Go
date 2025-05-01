package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	
)

var jwtkey = []byte("JWT_KEY")

func GenerateToken(email,role string) (string,error) {
	claims := jwt.MapClaims{
		"email":email,
		"role":role,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString(jwtkey)
}
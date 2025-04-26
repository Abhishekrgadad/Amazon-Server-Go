package middlewares

import (
	"os"
	"server/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func JWTProtected(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.UnauthorizedError(c, "Missing authentication token")
		}

		tokenString := strings.Split(authHeader, "Bearer ")
		if len(tokenString) != 2 {
			return errors.UnauthorizedError(c, "Invalid token format")
		}

		_ = godotenv.Load()
		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil || !token.Valid {
			return errors.UnauthorizedError(c, "Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return errors.UnauthorizedError(c, "Invalid token claims")
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			return errors.UnauthorizedError(c, "Role not found in token")
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Locals("user_id", claims["user_id"])
				c.Locals("role", userRole)
				return c.Next()
			}
		}

		return errors.ForbiddenError(c, "Access denied")
	}
}


// func JWTProtected(allowedRoles ...string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return errors.UnauthorizedError(c, "Missing authentication token")
// 		}

// 		tokenString := strings.Split(authHeader, "Bearer ")
// 		if len(tokenString) != 2 {
// 			return errors.UnauthorizedError(c, "Invalid token format")
// 		}

// 		godotenv.Load()
// 		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("JWT_KEY")), nil
// 		})

// 		if err != nil || !token.Valid {
// 			return errors.UnauthorizedError(c, "Invalid or expired token")
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			return errors.UnauthorizedError(c, "Invalid token claims")
// 		}

// 		userRole := claims["role"].(string)
// 		allowedRoles := []string{"seller", "admin"}

// 		for _, role := range allowedRoles {
// 			if userRole == role {
// 				c.Locals("user_id", claims["user_id"])
// 				c.Locals("role", userRole)
// 				return c.Next()
// 			}
// 		}

// 		return errors.ForbiddenError(c, "Access denied")
// 	}
// }

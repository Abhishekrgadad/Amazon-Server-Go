package auth

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"server/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Hash the password in bytes with defaultcost:10 (complexity)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Function compares the hashed string password with the user entered password for authentication
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Function generates the tokens after login successful
func GenerateJWT(userID, email, role string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_KEY")
	return token.SignedString([]byte(secret))
}

// Function generates reset tokens to reset the passwords.
func GenerateResetToken(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(10 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_KEY")
	return token.SignedString([]byte(secret))
}

// Function to fetch the respective collections by roles
func getCollectionByRole(role string) *mongo.Collection {
	switch role {
	case "customer":
		return config.DB.Collection("users")
	case "seller":
		return config.DB.Collection("sellers")
	case "admin":
		return config.DB.Collection("admins")
	case "user":
		return config.DB.Collection("users")
	default:
		return nil
	}
}

// Function to update with new password after reset
func ResetPassword(email, newPassword, role string) error {
	collection := getCollectionByRole(role)
	if collection == nil {
		return errors.New("invalid role")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"password": hashedPassword}})
	if err != nil {
		return err
	}
	return nil
}

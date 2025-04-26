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
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Function to fetch the user by email and password
func GetUsers(page int) ([]User, int64, int, error) {
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := 10
	skip := int64((page - 1) * limit)
	limitInt64 := int64(limit)
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, 0, err
	}
	// pagination
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limitInt64,
	})
	if err != nil {
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(totalCount) / limit
	if totalCount%int64(limit) != 0 {
		totalPages++
	}
	return users, totalCount, len(users), nil
}

// Function to fetch the seller by email and password
func GetSellers(page int) ([]Seller, int64, int, error) {
	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := 10
	skip := int64((page - 1) * limit)
	limitInt64 := int64(limit)
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, 0, err
	}
	// pagination
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limitInt64,
	})
	if err != nil {
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)

	var sellers []Seller
	if err = cursor.All(ctx, &sellers); err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(totalCount) / limit
	if totalCount%int64(limit) != 0 {
		totalPages++
	}
	return sellers, totalCount, len(sellers), nil
}

// Function to fetch the admin by email and password
func GetAdmins(page int) ([]Admin, int64, int, error) {
	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := 10
	skip := int64((page - 1) * limit)
	limitInt64 := int64(limit)
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, 0, err
	}
	// pagination
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limitInt64,
	})
	if err != nil {
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)

	var admins []Admin
	if err = cursor.All(ctx, &admins); err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(totalCount) / limit
	if totalCount%int64(limit) != 0 {
		totalPages++
	}
	return admins, totalCount, len(admins), nil
}

// Function to Register a new user
func RegisterUser(user User) (interface{}, error) {
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser User
	err := collection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"email": user.Email},
		{"phone_number": user.PhoneNumber},
	}}).Decode(&existingUser)
	if err == nil {
		return nil, mongo.ErrClientDisconnected // we'll handle error message in controller
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

// Function to Register a new seller
func RegisterSeller(seller Seller) (interface{}, error) {
	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingSeller Seller
	err := collection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"email": seller.Email},
		{"phone_number": seller.PhoneNumber},
	}}).Decode(&existingSeller)

	if err == nil {
		return nil, mongo.ErrClientDisconnected
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seller.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	seller.Password = string(hashedPassword)

	result, err := collection.InsertOne(ctx, seller)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

// Function to Register a new admin
func RegisterAdmin(admin Admin) (interface{}, error) {
	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingAdmin Admin
	err := collection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"email": admin.Email},
		{"phone_number": admin.PhoneNumber},
	}}).Decode(&existingAdmin)
	if err == nil {
		return nil, mongo.ErrClientDisconnected
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	admin.Password = string(hashedPassword)
	result, err := collection.InsertOne(ctx, admin)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

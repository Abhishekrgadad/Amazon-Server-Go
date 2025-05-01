package auth

import (
	"context"
	"errors"
	"fmt"
	"server/config"
	validation "server/modules/Validation"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func IsEmailTaken(email string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	checkmailcollection := []string{"users", "admins"}

	for _, col := range checkmailcollection {
		collection := config.DB.Collection(col)

		filter := bson.M{"email": email}

		var result User
		err := collection.FindOne(ctx, filter).Decode(&result)
		if err == nil {
			return true, col, nil
		}
		if err != mongo.ErrNoDocuments {
			return false, "", err
		}
	}
	return false, "", nil
}

func IsPhoneTaken(phone string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checkphonecollection := []string{"users", "admins"}

	for _, col := range checkphonecollection {
		collection := config.DB.Collection(col)
		filter := bson.M{"phone": phone}

		var result User
		err := collection.FindOne(ctx, filter).Decode(&result)
		if err == nil {
			return true, col, nil
		}
		if err != mongo.ErrNoDocuments {
			return false, "", err
		}
	}
	return false, "", nil
}

func RegisterUserHandler(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse Json"})
	}

	if err := validation.ValidateInputs(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	usercollection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checkMail, collectionName, err := IsEmailTaken(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking the existing user",
		})
	}
	if checkMail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Email already exists in '%s'", collectionName),
		})
	}

	checkphone, collectionName, err := IsPhoneTaken(user.Phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking the existng user",
		})
	}
	if checkphone {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Phone no. already exists in '%s'", collectionName),
		})
	}

	user.Password = HashPassword(user.Password)
	user.ID = primitive.NewObjectID()

	_, err = usercollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created successfully"})
}

func RegisterAdminHandler(c *fiber.Ctx) error {
	admin := new(Admin)
	if err := c.BodyParser(admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse admin"})
	}

	if err := validation.ValidateInputs(admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	admincollection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checkMail, collectionName, err := IsEmailTaken(admin.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking the existing user",
		})
	}
	if checkMail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Email already exists in '%s'", collectionName),
		})
	}

	checkphone, collectionName, err := IsPhoneTaken(admin.Phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking the existng user",
		})
	}
	if checkphone {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Phone no. already exists in '%s'", collectionName),
		})
	}

	admin.Password = HashPassword(admin.Password)
	admin.ID = primitive.NewObjectID()
	admin.Role = "admin"

	_, err = admincollection.InsertOne(ctx, admin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create admin user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "admin created successfully"})
}

func Login(input LoginRequest) (string,error) {
	var collectionName string

	switch input.Role {
	case "user": 
		collectionName = "users"
	case "admin":
		collectionName = "admins"
	default: 
		return "",errors.New("invalid role type")
	}

	collection := config.DB.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var dbUser User
	err := collection.FindOne(ctx,bson.M{"email":input.Email}).Decode(&dbUser)
	if err != nil {
		return "",errors.New("invalid email")
	}
	
	if !CheckPassword(input.Password,dbUser.Password) {
		return "", errors.New("invalid Password")
	} 
	
	token ,err := config.GenerateToken(dbUser.Email,dbUser.Role)
	if err != nil{
		return "",errors.New("failed to generate JWT ")
	}

	return token, nil
}
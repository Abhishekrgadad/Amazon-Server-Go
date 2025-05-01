package auth

import (
	"context"
	"server/config"
	validation "server/modules/Validation"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hash,_ := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	return string(hash)
}

func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed),[]byte(password))
	return err == nil
}

func RegisterUserHandler(c *fiber.Ctx) error{
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Cannot parse Json"})
	}

	if err := validation.ValidateUser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}

	usercollection := config.DB.Collection("users")
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	checkmail, err := usercollection.CountDocuments(ctx,bson.M{"email":user.Email})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Error checking the existing user"})
	}
	if checkmail > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"User already exist"})
	}
	checkphone, err := usercollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Error checking the existing phone no."})
	}
	if checkphone > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"user already exist"})
	}

	user.Password = HashPassword(user.Password)
	user.ID = primitive.NewObjectID()

	_,err = usercollection.InsertOne(ctx,user) 
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Failed to create user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message":"user created successfully"})
}

func RegisterAdmin(c *fiber.Ctx) error {
	admin := new(Admin)
	if err := c.BodyParser(admin); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Failed to parse admin"})
	}

	admincollection := config.DB.Collection("admins")
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	checkemail, err := admincollection.CountDocuments(ctx,bson.M{"email":admin.Email})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"error checking the existing email"})
	}
	if checkemail > 0{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Email already in use"})
	}
	checkphone,err := admincollection.CountDocuments(ctx,bson.M{"phone":admin.Phone})
	if err !=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"error checking the exising phone"})
	}
	if checkphone > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"phone already in use"})
	}

	admin.Password = HashPassword(admin.Password)
	admin.ID = primitive.NewObjectID()
	admin.Role = "admin"

	_,err = admincollection.InsertOne(ctx,admin)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Failed to create admin user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message":"admin created successfully"})
}
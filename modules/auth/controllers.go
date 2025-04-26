package auth

import (
	"context"
	"server/config"
	"server/errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// Function to register a user
func RegisterUserHandler(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	if err := validate.Struct(user); err != nil {
		return errors.BadRequestError(c, "Validation failed: "+err.Error())
	}
	id, err := RegisterUser(user)
	if err != nil {
		if err == mongo.ErrClientDisconnected {
			return errors.ConflictError(c, "Email or phone number already registered")
		}
		return errors.InternalServerError(c, "Failed to register user: "+err.Error())
	}
	return c.JSON(fiber.Map{"message": "User registered successfully", "user_id": id})
}

// Function to register a seller
func RegisterSellerHandler(c *fiber.Ctx) error {
	var seller Seller
	if err := c.BodyParser(&seller); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	if err := validate.Struct(seller); err != nil {
		return errors.BadRequestError(c, "Validation failed: "+err.Error())
	}
	id, err := RegisterSeller(seller)
	if err != nil {
		if err == mongo.ErrClientDisconnected {
			return errors.ConflictError(c, "Email or phone number already registered")
		}
		return errors.InternalServerError(c, "Failed to register seller: "+err.Error())
	}
	return c.JSON(fiber.Map{"message": "Seller registered successfully", "seller_id": id})
}

// Function to register an admin
func RegisterAdminHandler(c *fiber.Ctx) error {
	var admin Admin
	if err := c.BodyParser(&admin); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	if err := validate.Struct(admin); err != nil {
		return errors.BadRequestError(c, "Validation failed: "+err.Error())
	}
	id, err := RegisterAdmin(admin)
	if err != nil {
		if err == mongo.ErrClientDisconnected {
			return errors.ConflictError(c, "Email or phone number already registered")
		}
		return errors.InternalServerError(c, "Failed to register admin: "+err.Error())
	}
	return c.JSON(fiber.Map{"message": "Admin registered successfully", "admin_id": id})
}

// Function to login a user, seller, or admin
func LoginHandler(c *fiber.Ctx) error {
	var loginData LoginRequest
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.BodyParser(&loginData); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	var collection *mongo.Collection
	var user interface{}
	switch loginData.Role {
	case "user":
		collection = config.DB.Collection("users")
		user = &User{}
	case "seller":
		collection = config.DB.Collection("sellers")
		user = &Seller{}
	case "admin":
		collection = config.DB.Collection("admins")
		user = &Admin{}
	default:
		return errors.BadRequestError(c, "Invalid role type")
	}
	err := collection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.UnauthorizedError(c, "Invalid credentials")
		}
		return errors.InternalServerError(c, "Database error")
	}

	var hashedPassword string
	switch u := user.(type) {
	case *User:
		hashedPassword = u.Password
	case *Seller:
		hashedPassword = u.Password
	case *Admin:
		hashedPassword = u.Password
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginData.Password)); err != nil {
		return errors.UnauthorizedError(c, "Invalid credentials")
	}
	var userID string
	switch u := user.(type) {
	case *User:
		userID = u.ID.Hex()
	case *Seller:
		userID = u.ID.Hex()
	case *Admin:
		userID = u.ID.Hex()
	}
	token, err := config.GenerateJWT(userID, loginData.Email, loginData.Role)
	if err != nil {
		return errors.InternalServerError(c, "Could not generate token")
	}
	return c.JSON(fiber.Map{"message": "Login successful", "token": token})
}

// Function to Reset Password
func RequestPasswordReset(c *fiber.Ctx) error {
	var request ResetRequest
	if err := c.BodyParser(&request); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	collection := getCollectionByRole(request.Role)
	if collection == nil {
		return errors.BadRequestError(c, "Invalid role type")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := collection.FindOne(ctx, bson.M{"email": request.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.InternalServerError(c, "User not found")
		}
		return errors.InternalServerError(c, "Database error")
	}
	resetToken, err := HashPassword(request.Email + time.Now().String())
	if err != nil {
		return errors.InternalServerError(c, "Could not generate reset token")
	}
	_, err = collection.UpdateOne(ctx, bson.M{"email": request.Email}, bson.M{"$set": bson.M{"reset_token": resetToken}})
	if err != nil {
		return errors.InternalServerError(c, "Failed to generate reset token")
	}
	return c.JSON(fiber.Map{"message": "Password reset token generated", "reset_token": resetToken})
}

// Function to update password
func UpdatePassword(c *fiber.Ctx) error {
	var request UpdatePasswordRequest
	if err := c.BodyParser(&request); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	collection := getCollectionByRole(request.Role)
	if collection == nil {
		return errors.BadRequestError(c, "Invalid role type")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := collection.FindOne(ctx, bson.M{"email": request.Email, "reset_token": request.ResetToken}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.UnauthorizedError(c, "Invalid reset token or email")
		}
		return errors.InternalServerError(c, "Database error")
	}
	hashedPassword, err := HashPassword(request.NewPassword)
	if err != nil {
		return errors.InternalServerError(c, "Could not hash new password")
	}
	_, err = collection.UpdateOne(ctx, bson.M{"email": request.Email}, bson.M{"$set": bson.M{"password": hashedPassword}, "$unset": bson.M{"reset_token": ""}})
	if err != nil {
		return errors.InternalServerError(c, "Failed to update password")
	}
	return c.JSON(fiber.Map{"message": "Password updated successfully. Please log in again."})
}

// Function to get the collection based on role
func GetAllUsersHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return errors.BadRequestError(c, "Invalid page number")
	}
	users, totalCount, totalPages, err := GetUsers(page)
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch users")
	}
	return c.JSON(fiber.Map{
		"data":         users,
		"total_count":  totalCount,
		"current_page": page,
		"total_pages":  totalPages,
	})
}

// Function to get Users
func GetUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("users")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid user ID")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user bson.M
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return errors.NotFoundError(c, "User not found")
	}
	return c.JSON(user)
}

func UpdateUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("users")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid user ID")
	}
	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	delete(updateData, "email")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return errors.InternalServerError(c, "Failed to update user")
	}
	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

func DeleteUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid user ID")
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return errors.InternalServerError(c, "Failed to delete user")
	}
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

func GetAllSellerHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return errors.BadRequestError(c, "Invalid page number")
	}
	sellers, totalCount, totalPages, err := GetSellers(page)
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch sellers")
	}
	return c.JSON(fiber.Map{
		"data":         sellers,
		"total_count":  totalCount,
		"current_page": page,
		"total_pages":  totalPages,
	})
}

// Function to get Sellers
func GetSellerHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid sellers ID")
	}

	var user bson.M
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return errors.NotFoundError(c, "seller not found")
	}

	return c.JSON(user)
}

// Function to update Sellers
func UpdateSellerHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid seller ID")
	}

	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}

	delete(updateData, "email")
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return errors.InternalServerError(c, "Failed to update seller")
	}

	return c.JSON(fiber.Map{"message": "Seller updated successfully"})
}

// Function to Delete Sellers
func DeleteSellerHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid seller ID")
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return errors.InternalServerError(c, "Failed to delete seller")
	}
	return c.JSON(fiber.Map{"message": "Seller deleted successfully"})
}

// Function to get all admins
func GetAllAdminsHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return errors.BadRequestError(c, "Invalid page number")
	}

	admins, totalCount, totalPages, err := GetAdmins(page)
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch admins")
	}
	return c.JSON(fiber.Map{
		"data":         admins,
		"total_count":  totalCount,
		"current_page": page,
		"total_pages":  totalPages,
	})
}

// Function to get Admins
func GetAdminHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid admin ID")
	}
	var user bson.M
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return errors.NotFoundError(c, "admin not found")
	}
	return c.JSON(user)
}

// Function to update Admins
func UpdateAdminHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid admin ID")
	}

	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	delete(updateData, "email") // Prevent email updates directly
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return errors.InternalServerError(c, "Failed to update admin")
	}
	return c.JSON(fiber.Map{"message": "Admin updated successfully"})
}

// Function to Delete Admins
func DeleteAdminHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid admin ID")
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return errors.InternalServerError(c, "Failed to delete admin")
	}
	return c.JSON(fiber.Map{"message": "Admin deleted successfully"})
}

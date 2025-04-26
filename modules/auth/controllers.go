package auth

import (
	"context"
	"math"
	"net/http"
	"server/config"
	"server/errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

//This controller contains all the authentication part routes realted to users, sellers and admins.

func RegisterUserHandler(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}
	if err := validate.Struct(user); err != nil {
		return errors.BadRequestError(c, "Validation failed: "+err.Error())
	}

	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser User
	err := collection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"email": user.Email},
		{"phone_number": user.PhoneNumber},
	}}).Decode(&existingUser)
	if err == nil {
		return errors.ConflictError(c, "Email or phone number already registered")
	} else if err != mongo.ErrNoDocuments {
		return errors.InternalServerError(c, "Database error")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.InternalServerError(c, "Could not hash password")
	}
	user.Password = string(hashedPassword)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return errors.InternalServerError(c, "Failed to register user")
	}
	return c.JSON(fiber.Map{"message": "User registered successfully", "user_id": result.InsertedID})
}

func RegisterSellerHandler(c *fiber.Ctx) error {
	var seller Seller

	if err := c.BodyParser(&seller); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
	}

	if err := validate.Struct(seller); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seller.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	seller.Password = string(hashedPassword)

	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingSeller Seller
	err = collection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"email": seller.Email},
		{"phone_number": seller.PhoneNumber},
	}}).Decode(&existingSeller)

	if err == nil {
		return errors.ConflictError(c, "Email or phone number already registered")
	} else if err != mongo.ErrNoDocuments {
		return errors.InternalServerError(c, "Database error")
	}

	// Insert seller into MongoDB
	result, err := collection.InsertOne(ctx, seller)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register seller"})
	}

	return c.JSON(fiber.Map{"message": "Seller registered successfully", "seller_id": result.InsertedID})
}

func RegisterAdminHandler(c *fiber.Ctx) error {
	var admin Admin
	if err := c.BodyParser(&admin); err != nil {
		return errors.BadRequestError(c, "Invalid request data")
	}

	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.InternalServerError(c, "Could not hash password")
	}
	admin.Password = string(hashedPassword)

	var existingAdmin Admin
	err = collection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"email": admin.Email},
		{"phone_number": admin.PhoneNumber},
	}}).Decode(&existingAdmin)

	if err == nil {
		return errors.ConflictError(c, "Email or phone number already registered")
	} else if err != mongo.ErrNoDocuments {
		return errors.InternalServerError(c, "Database error")
	}

	result, err := collection.InsertOne(ctx, admin)
	if err != nil {
		return errors.InternalServerError(c, "Failed to register admin")
	}

	return c.JSON(fiber.Map{"message": "Admin registered successfully", "admin_id": result.InsertedID})
}

func LoginHandler(c *fiber.Ctx) error {
	var loginData LoginRequest
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func GetAllUsersHandler(c *fiber.Ctx) error {
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := strconv.Atoi(c.Params("page"))
	if err != nil || page < 1 {
		return errors.BadRequestError(c, "Invalid page number")
	}
	limit := 10
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return errors.InternalServerError(c, "Failed to count users")
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	if page > totalPages {
		return errors.NotFoundError(c, "No data available for this page")
	}

	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Limit: &limit64,
		Skip:  &skip,
	})
	if err != nil {
		return errors.InternalServerError(c, "Database error")
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		return errors.InternalServerError(c, "Failed to parse data")
	}

	return c.JSON(fiber.Map{
		"data":        users,
		"total_users": totalCount,
		"page_no":     page,
		"total_pages": totalPages,
	})
}

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
	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := strconv.Atoi(c.Params("page"))
	if err != nil || page < 1 {
		return errors.BadRequestError(c, "Invalid page number")
	}

	limit := 10
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return errors.InternalServerError(c, "Failed to count sellers")
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	if page > totalPages {
		return errors.NotFoundError(c, "No data available for this page")
	}
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Limit: &limit64,
		Skip:  &skip,
	})
	if err != nil {
		return errors.InternalServerError(c, "Database error")
	}
	defer cursor.Close(ctx)

	var seller []bson.M
	if err = cursor.All(ctx, &seller); err != nil {
		return errors.InternalServerError(c, "Failed to parse data")
	}

	return c.JSON(fiber.Map{
		"data":        seller,
		"total_count": totalCount,
		"page_no":     page,
		"total_pages": totalPages,
	})
}

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

func DeleteSellerHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("sellers")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid seller ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return errors.InternalServerError(c, "Failed to delete seller")
	}

	return c.JSON(fiber.Map{"message": "Seller deleted successfully"})
}

func GetAllAdminsHandler(c *fiber.Ctx) error {
	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := strconv.Atoi(c.Params("page"))
	if err != nil || page < 1 {
		return errors.BadRequestError(c, "Invalid page number")
	}

	limit := 10
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return errors.InternalServerError(c, "Failed to count admins")
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	if page > totalPages {
		return errors.NotFoundError(c, "No data available for this page")
	}
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Limit: &limit64,
		Skip:  &skip,
	})
	if err != nil {
		return errors.InternalServerError(c, "Database error")
	}
	defer cursor.Close(ctx)

	var admins []bson.M
	if err = cursor.All(ctx, &admins); err != nil {
		return errors.InternalServerError(c, "Failed to parse data")
	}

	return c.JSON(fiber.Map{
		"data":         admins,
		"total_admins": totalCount,
		"page_no":      page,
		"total_pages":  totalPages,
	})
}

func GetAdminHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.DB.Collection("admins")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BadRequestError(c, "Invalid admin ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user bson.M
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return errors.NotFoundError(c, "admin not found")
	}

	return c.JSON(user)
}

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

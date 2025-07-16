package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"server/config"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// func GenerateJWT(userID, email, role string) (string, error) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	claims := jwt.MapClaims{
// 		"user_id": userID,
// 		"email":   email,
// 		"role":    role,
// 		"exp":     time.Now().Add(time.Hour * 24).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	secret := os.Getenv("JWT_KEY")
// 	return token.SignedString([]byte(secret))
// }

func IsEmailTaken(email string) (bool, error) {
	collections := []string{"users", "sellers", "admins"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, colName := range collections {
		collection := config.DB.Collection(colName)

		if email != "" {
			err := collection.FindOne(ctx, bson.M{"email": email}).Err()
			if err == nil {
				return true, nil
			} else if err != mongo.ErrNoDocuments {
				return false, err
			}
		}
	}
	return false, nil
}

func IsPhoneTaken(phone_number string) (bool, error) {
	collections := []string{"users", "sellers", "admins"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, colName := range collections {
		collection := config.DB.Collection(colName)

		if phone_number != "" {
			err := collection.FindOne(ctx, bson.M{"phone_number": phone_number}).Err()
			if err == nil {
				return true, nil
			} else if err != mongo.ErrNoDocuments {
				return false, err
			}
		}
	}
	return false, nil
}

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

func ResetPassword(email, newPassword, role string) error {
	collection := getCollectionByRole(role)
	if collection == nil {
		return errors.New("invalid role")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func GetUsers(page int) ([]User, int64, int, error) {
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	limit := int64(10)
	skip := int64(page-1) * limit
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, 0, err
	}
	project := bson.M{
		"password": 0, 
	}
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetProjection(project).SetSort(bson.D{{Key: "fullname",Value: 1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, 0, err
	}
	return users, totalCount, len(users), nil
}

func GetSellers(page int) ([]Seller, int64, int, error) {
	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	limit := int64(10)
	skip := int64(page-1) * limit
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, 0, err
	}
	project := bson.M{
		"password": 0,
	}
	filter := options.Find().SetSkip(skip).SetLimit(limit).SetProjection(project)
	cursor, err := collection.Find(ctx, bson.M{}, filter)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)

	var sellers []Seller
	if err = cursor.All(ctx, &sellers); err != nil {
		return nil, 0, 0, err
	}
	return sellers, totalCount, len(sellers), nil
}

func GetAdmins(page int) ([]Admin, int64, int, error) {
	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	limit := int64(10)
	skip := int64(page-1) * limit
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, 0, err
	}
	project := bson.M{
		"password": 0,
	}
	filter := options.Find().SetSkip(skip).SetLimit(limit).SetProjection(project)
	cursor, err := collection.Find(ctx, bson.M{},filter)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)

	var admins []Admin
	if err = cursor.All(ctx, &admins); err != nil {
		return nil, 0, 0, err
	}
	return admins, totalCount, len(admins), nil
}

func RegisterUser(user User) (interface{}, error) {
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	emailExist, err := IsEmailTaken(user.Email)
	if emailExist {
		return nil, errors.New("email already in use")
	}
	if err != nil {
		return nil, err
	}
	phoneExist, err := IsPhoneTaken(user.PhoneNumber)
	if phoneExist {
		return nil, errors.New("phone already in use")
	}
	if err != nil {
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

func RegisterSeller(seller Seller) (interface{}, error) {
	collection := config.DB.Collection("sellers")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	emailExist, err := IsEmailTaken(seller.Email)
	if emailExist {
		return nil, errors.New("email already in use")
	}
	if err != nil {
		return nil, err
	}
	phoneExist, err := IsPhoneTaken(seller.PhoneNumber)
	if phoneExist {
		return nil, errors.New("phone already in use")
	}
	if err != nil {
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

func RegisterAdmin(admin Admin) (interface{}, error) {
	collection := config.DB.Collection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	emailExist, err := IsEmailTaken(admin.Email)
	if emailExist {
		return nil, errors.New("email already in use")
	}
	if err != nil {
		return nil, err
	}
	phoneExist, err := IsPhoneTaken(admin.PhoneNumber)
	if phoneExist {
		return nil, errors.New("phone already in use")
	}
	if err != nil {
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

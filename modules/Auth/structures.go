package auth

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email" `
	Phone    string             `json:"phone" bson:"phone" validate:"required,e164,len=13"`
	Password string             `json:"password" bson:"password,omitempty" validate:"required,min=6"`
	Role     string             `json:"role" bson:"role" validate:"required,oneof=user admin"`
}

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email" `
	Phone    string             `json:"phone" bson:"phone" validate:"required,e164,len=13"`
	Password string             `json:"password" bson:"password,omitempty" validate:"required,min=6"`
	Role     string             `json:"role" bson:"role" validate:"required,oneof=user admin"`
}

var RegisterRequest struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Phone    string             `json:"phone" bson:"phone" validate:"required,e164,len=13"`
	Password string             `json:"password" bson:"password,omitempty" validate:"required,min=6"`
	Role     string             `json:"role" bson:"role" validate:"required,oneof=user admin"`
}

var LoginRequest struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password,omitempty" validate:"required,min=6"`
	Role     string `json:"role" bson:"role" validate:"required,oneof=user admin"`
}

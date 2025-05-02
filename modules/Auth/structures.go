package auth

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name" validate:"required"`
	Email    string             `bson:"email" json:"email" validate:"required,email" `
	Phone    string             `bson:"phone" json:"phone" validate:"required,e164,len=13"`
	Password string             `bson:"password,omitempty" json:"password" validate:"required,min=6"`
	Role     string             `bson:"role" json:"role" validate:"required,oneof=user admin"`
}

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name" validate:"required"`
	Email    string             `bson:"email" json:"email" validate:"required,email" `
	Phone    string             `bson:"phone" json:"phone" validate:"required,e164,len=13"`
	Password string             `bson:"password,omitempty" json:"password" validate:"required,min=6"`
	Role     string             `bson:"role" json:"role" validate:"required,oneof=user admin"`
}

type RegisterRequest struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name" validate:"required"`
	Email    string             `bson:"email" json:"email" validate:"required,email"`
	Phone    string             `bson:"phone" json:"phone" validate:"required,e164,len=13"`
	Password string             `bson:"password,omitempty" json:"password" validate:"required,min=6"`
	Role     string             `bson:"role" json:"role" validate:"required,oneof=user admin"`
}

type LoginRequest struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password,omitempty" json:"-" validate:"required,min=6"`
	Role     string `bson:"role" json:"role" validate:"required,oneof=user admin"`
}

type Response struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email" `
	Phone    string             `bson:"phone" json:"phone" `
	Password string             `bson:"password,omitempty" json:"-" `
	Role     string             `bson:"role" json:"role"`
}
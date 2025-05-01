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
	Name     string             `bson:"name" json:"name" validate:"required"`
	Email    string             `bson:"email" json:"email" validate:"required,email"`
	Phone    string             `bson:"phone" json:"phone" validate:"required,e164,len=13"`
	Password string             `bson:"password,omitempty" json:"password" validate:"required,min=6"`
	Role     string             `bson:"role" json:"role" validate:"required"`
}

var RegisterRequest struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required,email" `
	Phone    string             `json:"phone" bson:"phone" validate:"required,e164,len=13"`
	Password string             `json:"password" bson:"password,omitempty" validate:"required,min=6"`
	Role     string             `json:"role" bson:"role" validate:"required,oneof=user admin"`
}

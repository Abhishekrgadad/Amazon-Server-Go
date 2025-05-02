package product

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Price       float64            `bson:"price" json:"price" validate:"required,gt=0"`
	Category    string             `bson:"category" json:"category" validate:"required"`
	Stock       int                `bson:"stock" json:"stock" validate:"required,gte=0"`
	CreatedAt   primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
}


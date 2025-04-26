package review

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"order_id"`
	Rating    int                `bson:"rating" json:"rating"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type ReviewRequest struct {
	ProductID string `json:"product_id"`
	OrderID   string `json:"order_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}
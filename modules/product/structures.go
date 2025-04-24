package product

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product Structure
type Product struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name" validate:"required"`
	Description    string             `json:"description" bson:"description" validate:"required"`
	Price          float64            `json:"price" bson:"price" validate:"required"`
	StockQuantity  int                `json:"stock_quantity" bson:"stock_quantity" validate:"required"`
	Category       string             `json:"category" bson:"category" validate:"required"`
	Brand          string             `json:"brand" bson:"brand" validate:"required"`
	Visibility     bool               `json:"visibility" bson:"visibility" validate:"required"`
	AverageRating  float64            `bson:"average_rating,omitempty" json:"average_rating,omitempty"`
	TotalReviews   int                `bson:"total_reviews"`
	ReviewComments []string           `bson:"review_comments,omitempty" json:"review_comments,omitempty"`
	CreatedAt      string             `json:"created_at" bson:"created_at"`
	UpdatedAt      string             `json:"updated_at" bson:"updated_at"`
}

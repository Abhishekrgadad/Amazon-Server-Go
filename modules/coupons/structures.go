package coupons

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code       string             `bson:"code" json:"code"`
	Discount   float64            `bson:"discount" json:"discount"`
	Active     bool               `bson:"active" json:"active"`
	IsPercent  bool               `bson:"is_percent" json:"is_percent"`
	ExpiryDate time.Time          `bson:"expiry_date" json:"expiry_date"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

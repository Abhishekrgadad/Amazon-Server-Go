package order

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckoutRequest struct {
	UserID      string `json:"user_id"`
	CartID      string `json:"cart_id"`
	PaymentType string `json:"payment_type"`
	Address     string `json:"address"`
	CouponCode  string `json:"coupon_code"`
}

type CartItem struct {
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductName string             `json:"product_name"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Price       float64            `bson:"price" json:"price"`
}

type CartItemResponse struct {
	ProductID   primitive.ObjectID `json:"product_id"`
	ProductName string             `json:"product_name"`
	Price       float64            `json:"price"`
	Quantity    int                `json:"quantity"`
	SubTotal    float64            `json:"sub_total"`
}

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Address     string             `bson:"address" json:"address"`
	Items       []CartItem         `bson:"items" json:"items"`
	TotalPrice  float64            `bson:"total_price" json:"total_price"`
	PaymentType string             `bson:"payment_type" json:"payment_type"`
	CouponCode  string             `bson:"coupon_code" json:"coupon_code"`
	Discount    float64            `bson:"discount" json:"discount"`
	Status      string             `bson:"status" json:"status"`
	UserDetails User               `bson:"user_details" json:"user_details"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type OrderNotification struct {
	OrderID    string      `json:"order_id"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	Items      interface{} `json:"items"`
	Address    string      `json:"address"`
	CreatedAt  string      `json:"created_at"`
}

type User struct {
	FullName        string `json:"full_name" bson:"full_name" validate:"required,min=3,max=20"`
	Email           string `json:"email" bson:"email" validate:"required,email"`
	PhoneNumber     string `json:"phone_number" bson:"phone_number" validate:"required,e164,len=13"`
	ShippingAddress string `json:"shipping_address" bson:"shipping_address"`
}

type StatusUpdate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Address     string             `bson:"address" json:"address"`
	Items       interface{}       `bson:"items" json:"items"`
	TotalPrice  float64            `bson:"total_price" json:"total_price"`
	PaymentType string             `bson:"payment_type" json:"payment_type"`
	CouponCode  string             `bson:"coupon_code" json:"coupon_code"`
	Discount    float64            `bson:"discount" json:"discount"`
	Status      string             `bson:"status" json:"status"`
	UserDetails interface{}               `bson:"user_details" json:"user_details"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

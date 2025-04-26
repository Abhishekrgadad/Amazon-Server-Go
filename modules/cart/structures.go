package cart

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddToCartItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type AddToCartRequest struct {
	UserID string          `json:"user_id"`
	Items  []AddToCartItem `json:"items"`
}

type AddCartResponse struct {
	CartID  primitive.ObjectID `json:"cart_id"`
	Details []CartItemDetail   `json:"items"`
	Total   float64            `json:"total_price"`
}

type UpdateCartRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CartItem struct {
	ProductID primitive.ObjectID `bson:"product_id"`
	Quantity  int                `bson:"quantity"`
	Price     float64            `bson:"price" json:"price"`
}

type Cart struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id"`
	Items  []CartItem         `bson:"items"`
}
type CartItemDetail struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
}

type CartItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type CartResponse struct {
	UserID     string             `json:"user_id"`
	CartID     string             `json:"cart_id"`
	Items      []CartItemResponse `json:"items"`
	TotalPrice float64            `json:"total_price"`
	Message    string             `json:"message"`
}

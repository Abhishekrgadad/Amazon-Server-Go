package order

import (
	"context"
	"fmt"
	"server/config"
	"server/modules/cart"
	"server/modules/coupons"
	"server/modules/product"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Function for placing order with the product details and price
func PlaceOrder(userID, cartID primitive.ObjectID, paymentType, address, couponCode string) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	cartdata, err := GetCartByID(cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %v", err)
	}
	if len(cartdata.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}
	totalPrice := 0.0

	prodcollection := config.DB.Collection("products")
	var updatedItems []CartItemResponse
	for _, item := range cartdata.Items {
		var product product.Product
		err := prodcollection.FindOne(ctx, bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product details: %v", item.ProductID.Hex())
		}
		if product.Price == 0 {
			return nil, fmt.Errorf("product %v has an invalid price", item.ProductID.Hex())
		}
		subTotal := product.Price * float64(item.Quantity)
		totalPrice += subTotal
		updatedItem := CartItemResponse{
			ProductID:   item.ProductID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			SubTotal:    subTotal,
		}
		updatedItems = append(updatedItems, updatedItem)
	}

	user, err := GetUserDetails(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user details: %v", err)
	}
	orderAddress := address
	if address == "" {
		orderAddress = user.ShippingAddress
	}

	var discount float64
	var appliedCoupon string
	couponcollection := config.DB.Collection("coupons")
	if couponCode != "" {
		var coupon coupons.Coupon
		err := couponcollection.FindOne(ctx, bson.M{
			"code":        couponCode,
			"active":      true,
			"expiry_date": bson.M{"$gt": time.Now()},
		}).Decode(&coupon)

		if err != nil {
			return nil, fmt.Errorf("invalid or expired coupon")
		}

		if !coupon.Active || time.Now().After(coupon.ExpiryDate) {
			return nil, fmt.Errorf("coupon is inactive or expired")
		}
		appliedCoupon = coupon.Code
		if coupon.IsPercent {
			discount = totalPrice * (coupon.Discount / 100)
		} else {
			discount = coupon.Discount
		}
	}

	finalTotal := totalPrice - discount
	if finalTotal < 0 {
		finalTotal = 0
	}

	var cartItems []CartItem
	for _, item := range updatedItems {
		cartItems = append(cartItems, CartItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})

		var product product.Product
		err := prodcollection.FindOne(ctx, bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product details for stock update: %v", item.ProductID.Hex())
		}

		if product.StockQuantity < item.Quantity {
			return nil, fmt.Errorf("not enough stock for product %v", item.ProductID.Hex())
		}

		StockQuantity := product.StockQuantity - item.Quantity
		_, err = prodcollection.UpdateOne(ctx, bson.M{"_id": item.ProductID}, bson.M{	
			"$set": bson.M{"stock_quantity": StockQuantity},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update stock for product %v: %v", item.ProductID.Hex(), err)
		}
		// err = prodcollection.FindOne(ctx, bson.M{"_id": item.ProductID}).Decode(&product)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to fetch product details for stock update: %v", item.ProductID.Hex())
		// }
		if product.StockQuantity == 0 {
			_, err := prodcollection.UpdateOne(ctx, bson.M{"_id": item.ProductID}, bson.M{
				"$set": bson.M{"visibility": false},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update stock for product %v: %v", item.ProductID.Hex(), err)
			}
		}

	}
	order := &Order{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Items:       cartItems,
		TotalPrice:  finalTotal,
		PaymentType: paymentType,
		Address:     orderAddress,
		Status:      "Order Confirmed",
		CreatedAt:   time.Now(),
		UserDetails: user,
		CouponCode:  appliedCoupon,
		Discount:    discount,
	}
	_, err = config.DB.Collection("orders").InsertOne(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %v", err)
	}
	return order, nil
}

// Function for getting cart details by ID
func GetCartByID(cartID primitive.ObjectID) (*cart.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart cart.Cart
	err := config.DB.Collection("cart").FindOne(ctx, bson.M{"_id": cartID}).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// Function for getting user details whos order is placed
func GetUserDetails(userID primitive.ObjectID) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := config.DB.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	return user, err
}

// Function for viewing all orders placed by the user
func ViewAllOrders(userID primitive.ObjectID, page int) ([]Order, int64, int64, error) {
	collection := config.DB.Collection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := int64(10)
	skip := (int64(page) - 1) * limit
	filter := bson.M{"user_id": userID}
	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count orders: %v", err)
	}
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"name": 1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch orders: %v", err)
	}
	defer cursor.Close(ctx)

	var orders []Order
	for cursor.Next(ctx) {
		var order Order
		if err := cursor.Decode(&order); err != nil {
			return nil, 0, 0, err
		}
		orders = append(orders, order)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, 0, err
	}
	totalPages := (totalCount + limit - 1) / limit
	return orders, totalCount, totalPages, nil
}

// Function for Canceling an order
func CancelOrder(orderID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": orderID, "user_id": userID}
	update := bson.M{"$set": bson.M{"status": "Canceled"}}
	res, err := config.DB.Collection("orders").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %v", err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("order not found or already canceled")
	}
	return nil
}

// Function for returning an order
func ReturnOrder(orderID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order Order
	err := config.DB.Collection("orders").FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	if order.Status != "Order Confirmed" && order.Status != "Delivered" {
		return fmt.Errorf("order cannot be returned or return request is already done")
	}
	_, err = config.DB.Collection("orders").UpdateOne(ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"status": "Returned"}})
	if err != nil {
		return fmt.Errorf("failed to update return status: %v", err)
	}
	return nil
}

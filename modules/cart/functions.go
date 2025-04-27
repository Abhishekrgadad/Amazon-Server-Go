package cart

import (
	"context"
	"errors"
	"fmt"
	"server/config"
	"server/modules/product"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Function to add products to Cart
func AddToCart(userID primitive.ObjectID, items []CartItem) (*AddCartResponse, error) {
	cartCollection := config.DB.Collection("cart")
	productCollection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart Cart
	err := cartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		newCart := Cart{
			ID:     primitive.NewObjectID(),
			UserID: userID,
			Items:  items,
		}
		_, err := cartCollection.InsertOne(ctx, newCart)
		if err != nil {
			return nil, fmt.Errorf("failed to create cart")
		}
		cart = newCart
	} else if err == nil {
		for _, newItem := range items {
			found := false
			for i, existingItem := range cart.Items {
				if existingItem.ProductID == newItem.ProductID {
					cart.Items[i].Quantity += newItem.Quantity
					found = true
					break
				}
			}
			if !found {
				cart.Items = append(cart.Items, newItem)
			}
		}
		_, err = cartCollection.UpdateOne(ctx,
			bson.M{"_id": cart.ID},
			bson.M{"$set": bson.M{"items": cart.Items}},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update cart")
		}
	} else {
		return nil, err
	}

	var details []CartItemDetail
	var total float64
	for _, item := range cart.Items {
		var product struct {
			Name  string  `bson:"name"`
			Price float64 `bson:"price"`
		}
		err := productCollection.FindOne(ctx, bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			continue
		}
		subtotal := float64(item.Quantity) * product.Price
		total += subtotal

		details = append(details, CartItemDetail{
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Price:       product.Price,
			Total:       subtotal,
		})
	}
	return &AddCartResponse{
		CartID:  cart.ID,
		Details: details,
		Total:   total,
	}, nil
}

// Function to view the cart items
func GetCart(userID string) (*CartResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cartCollection := config.DB.Collection("cart")
	productCollection := config.DB.Collection("products")
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var cart Cart
	err = cartCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("Cart is empty")
	} else if err != nil {
		return nil, errors.New("failed to retrieve cart")
	}

	var responseItems []CartItemResponse
	var total float64
	for _, item := range cart.Items {
		var product product.Product
		err := productCollection.FindOne(ctx, bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			continue // Skip if product not found
		}
		subtotal := product.Price * float64(item.Quantity)
		total += subtotal
		responseItems = append(responseItems, CartItemResponse{
			ProductID:   product.ID.Hex(),
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}
	return &CartResponse{
		UserID:     cart.UserID.Hex(),
		CartID:     cart.ID.Hex(),
		Items:      responseItems,
		TotalPrice: total,
		Message:    "Cart details fetched successfully",
	}, nil
}

// Function to update Cart details
func UpdateCartItem(userID string, productID string, quantity int) error {
	collection := config.DB.Collection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	productObjID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New("invalid product ID format")
	}
	filter := bson.M{"user_id": userObjID, "items.product_id": productObjID}
	update := bson.M{"$set": bson.M{"items.$.quantity": quantity}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %v", err)
	}
	if result.MatchedCount == 0 {
		pushUpdate := bson.M{"$push": bson.M{"items": bson.M{"product_id": productObjID, "quantity": quantity}}}
		_, err := collection.UpdateOne(ctx, bson.M{"user_id": userObjID}, pushUpdate)
		if err != nil {
			return fmt.Errorf("product not in cart, and failed to add it: %v", err)
		}
	}
	return nil
}

// Functino to Remove items from cart
func RemoveCartItem(userID string, productID string) error {
	collection := config.DB.Collection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}
	productObjID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New("invalid product ID")
	}
	_, err = collection.UpdateOne(ctx, bson.M{"user_id": userObjID},
		bson.M{"$pull": bson.M{"items": bson.M{"product_id": productObjID}}})
	return err
}

// Function to clear full cart from DB
func ClearCart(userID string) error {
	collection := config.DB.Collection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}
	_, err = collection.DeleteOne(ctx, bson.M{"user_id": userObjID})
	return err
}

// Function to fetch all cart records
func GetAllCarts() ([]Cart, error) {
	cartCollection := config.DB.Collection("cart")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var carts []Cart
	cursor, err := cartCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch carts: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var cart Cart
		if err := cursor.Decode(&cart); err != nil {
			return nil, fmt.Errorf("failed to decode cart: %v", err)
		}
		carts = append(carts, cart)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}
	return carts, nil
}

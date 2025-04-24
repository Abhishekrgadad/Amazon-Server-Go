package coupons

import (
	"context"
	"time"

	"server/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCoupon(coupon *Coupon) error {
	coupon.ID = primitive.NewObjectID()
	coupon.CreatedAt = time.Now()
	_, err := config.DB.Collection("coupons").InsertOne(context.TODO(), coupon)
	return err
}

func GetAllCoupons() ([]Coupon, error) {
	var coupons []Coupon
	cursor, err := config.DB.Collection("coupons").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var c Coupon
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		coupons = append(coupons, c)
	}
	return coupons, nil
}

func GetCouponByID(id primitive.ObjectID) (*Coupon, error) {
	var coupon Coupon
	err := config.DB.Collection("coupons").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&coupon)
	return &coupon, err
}

func UpdateCoupon(id primitive.ObjectID, updated bson.M) error {
	_, err := config.DB.Collection("coupons").UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updated})
	return err
}

func DeleteCoupon(id primitive.ObjectID) error {
	_, err := config.DB.Collection("coupons").DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

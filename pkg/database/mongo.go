package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect thiết lập kết nối đến MongoDB với URI được cung cấp.
func Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	// Kiểm tra kết nối
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

// Collection alias để dễ sử dụng.
type Collection = *mongo.Collection

// InsertNotification chèn thông báo vào collection.
func InsertNotification(coll Collection, notification interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.InsertOne(ctx, notification)
	if err != nil {
		log.Printf("Lỗi chèn thông báo: %v", err)
	}
	return err
}

// UpdateNotificationStatus cập nhật trạng thái thông báo dựa trên ID.
func UpdateNotificationStatus(coll Collection, id string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Lỗi cập nhật trạng thái cho notification %s: %v", id, err)
	}
	return err
}

package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabase implements Database interface
type MongoDatabase struct {
	db *mongo.Database
}

func Connect(ctx context.Context) (*MongoDatabase, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	db := client.Database("rakuten_auth")

	// Create unique index on email to prevents duplicate accounts
	_, err = db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return &MongoDatabase{db: db}, nil
}

func (m *MongoDatabase) CreateUser(ctx context.Context, user User) error {
	_, err := m.db.Collection("users").InsertOne(ctx, user)
	return err
}

func (m *MongoDatabase) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := m.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *MongoDatabase) IsDuplicateKeyError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

func (m *MongoDatabase) Disconnect(ctx context.Context) error {
	return m.db.Client().Disconnect(ctx)
}

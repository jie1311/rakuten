package database

import "context"

type User struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"-"`
}

// Database interface for dependency injection and testing
type Database interface {
	CreateUser(ctx context.Context, user User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	IsDuplicateKeyError(err error) bool
	Disconnect(ctx context.Context) error
}

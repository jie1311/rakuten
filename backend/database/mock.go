package database

import (
	"context"
	"errors"
)

// MockDatabase implements Database interface for testing
type MockDatabase struct {
	Users map[string]User
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		Users: make(map[string]User),
	}
}

func (m *MockDatabase) CreateUser(ctx context.Context, user User) error {
	if _, exists := m.Users[user.Email]; exists {
		return errors.New("duplicate key error")
	}
	m.Users[user.Email] = user
	return nil
}

func (m *MockDatabase) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	user, exists := m.Users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (m *MockDatabase) IsDuplicateKeyError(err error) bool {
	return err != nil && err.Error() == "duplicate key error"
}

func (m *MockDatabase) Disconnect(ctx context.Context) error {
	return nil
}

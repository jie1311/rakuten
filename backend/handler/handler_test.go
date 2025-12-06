package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jie1311/rakuten/backend/database"
)

func TestHandleSignup(t *testing.T) {
	// Create mock database
	mockDB := database.NewMockDatabase()
	handler := NewHandler(mockDB)

	// Create request
	reqBody := SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.HandleSignup(w, req)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Verify user was created in mock DB
	user, err := mockDB.FindUserByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Errorf("User was not created: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}
}

func TestHandleSignup_DuplicateEmail(t *testing.T) {
	mockDB := database.NewMockDatabase()
	handler := NewHandler(mockDB)

	// Pre-populate user
	mockDB.Users["test@example.com"] = database.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	reqBody := SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleSignup(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
}

func TestHandleSignin(t *testing.T) {
	mockDB := database.NewMockDatabase()
	handler := NewHandler(mockDB)

	// Create a user first
	reqBody := SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.HandleSignup(w, req)

	// Now test signin
	signinBody := SigninRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ = json.Marshal(signinBody)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	handler.HandleSignin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify response contains token
	var response AuthResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.Token == "" {
		t.Error("Expected token in response")
	}
	if response.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", response.Email)
	}
}

func TestHandleSignin_InvalidCredentials(t *testing.T) {
	mockDB := database.NewMockDatabase()
	handler := NewHandler(mockDB)

	signinBody := SigninRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(signinBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleSignin(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestHandleMe(t *testing.T) {
	mockDB := database.NewMockDatabase()
	handler := NewHandler(mockDB)

	// Create a user and get token
	reqBody := SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.HandleSignup(w, req)

	signinBody := SigninRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ = json.Marshal(signinBody)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	handler.HandleSignin(w, req)

	var authResponse AuthResponse
	json.NewDecoder(w.Body).Decode(&authResponse)

	// Test /me endpoint with token
	req = httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+authResponse.Token)
	w = httptest.NewRecorder()

	// Wrap with auth middleware
	handler.AuthMiddleware(http.HandlerFunc(handler.HandleMe)).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var meResponse MeResponse
	json.NewDecoder(w.Body).Decode(&meResponse)
	if meResponse.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", meResponse.Email)
	}
}

func TestHandleSignout(t *testing.T) {
	mockDB := database.NewMockDatabase()
	handler := NewHandler(mockDB)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signout", nil)
	w := httptest.NewRecorder()

	handler.HandleSignout(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockDB := database.NewMockDatabase()
	handler := NewHandler(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()

	handler.AuthMiddleware(http.HandlerFunc(handler.HandleMe)).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

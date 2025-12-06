package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jie1311/rakuten/backend/database"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("rakuten_secret_key")

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type MeResponse struct {
	Email string `json:"email"`
}

type Handler struct {
	DB database.Database
}

func NewHandler(db database.Database) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req, ok := decodeJSON[SignupRequest](w, r)
	if !ok {
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = h.DB.CreateUser(ctx,
		database.User{
			Email:    req.Email,
			Password: string(hashedPassword),
		})

	if err != nil {
		if h.DB.IsDuplicateKeyError(err) {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) HandleSignin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req, ok := decodeJSON[SigninRequest](w, r)
	if !ok {
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	user, err := h.DB.FindUserByEmail(ctx, req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(req.Email)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	encodeJSON(w, AuthResponse{Token: token, Email: req.Email})
}

func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	encodeJSON(w, MeResponse{Email: email})
}

func (h *Handler) HandleSignout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(r.Context(), "email", claims["email"].(string))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return data, false
	}
	return data, true
}

func encodeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func generateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

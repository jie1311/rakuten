package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jie1311/rakuten/backend/database"
	"github.com/jie1311/rakuten/backend/handler"
)

func main() {
	// connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect(ctx)

	// create handler with database dependency
	h := handler.NewHandler(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // allow frontend
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/signup", h.HandleSignup)
		r.Post("/signin", h.HandleSignin)
		r.Post("/signout", h.HandleSignout)
	})

	r.With(h.AuthMiddleware).Get("/api/me", h.HandleMe) // protected route

	log.Println("Server starting on 0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", r)
}

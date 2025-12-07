.PHONY: help install dev build clean backend frontend test mongodb

help:
	@echo "Available commands:"
	@echo "  make install   - Install dependencies for both frontend and backend"
	@echo "  make dev       - Run MongoDB, frontend and backend in development mode"
	@echo "  make mongodb   - Start MongoDB server"
	@echo "  make backend   - Run backend server only"
	@echo "  make frontend  - Run frontend dev server only"
	@echo "  make build     - Build both frontend and backend"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make test      - Run tests"

install:
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

dev:
	@echo "Starting MongoDB, frontend and backend..."
	@make -j3 mongodb backend frontend

mongodb:
	@echo "Starting MongoDB server..."
	mongod --dbpath ~/mongodb/data --bind_ip 0.0.0.0 --port 27017

backend:
	@echo "Starting backend server..."
	cd backend && go run main.go

frontend:
	@echo "Starting frontend dev server..."
	cd frontend && npm run dev

build:
	@echo "Building backend..."
	cd backend && go build -o backend main.go
	@echo "Building frontend..."
	cd frontend && npm run build

clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/backend
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	@echo "Clean complete"

test:
	@echo "Running backend tests..."
	cd backend && go test ./...
	@echo "Running frontend tests..."
	cd frontend && npm test

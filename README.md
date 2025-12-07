# Rakuten

A full-stack web application with Go backend and React frontend, featuring user authentication.

## Prerequisites

- Go 1.24+
- npm or yarn
- MongoDB

### Setup MongoDB Data Directory

Create a directory for MongoDB data:
```bash
mkdir -p ~/mongodb/data
```

## Quick Start

### 1. Clone and Install

```bash
# Install all dependencies
make install
```

### 2. Configure Environment

Create environment files:

**Backend** (`/backend/.env`):
```env
JWT_SECRET=your-secret-key-change-this-in-production
PORT=8080
```

**Frontend** (`/frontend/.env`):
```env
VITE_API_URL=http://localhost:8080
```

### 3. Run the Application

```bash
# Run both frontend and backend concurrently
make dev
```

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:8081`

### Run Separately

```bash
# Backend only
make backend

# Frontend only
make frontend
```

## Building for Production

```bash
# Build both frontend and backend
make build
```

Output:
- Backend binary: `/backend/backend`
- Frontend static files: `/frontend/dist/`

## Available Make Commands

```bash
make help      # Show all available commands
make install   # Install dependencies
make dev       # Run in development mode
make backend   # Run backend only
make frontend  # Run frontend only
make build     # Build for production
make clean     # Remove build artifacts
make test      # Run tests
```

## Features

- ✅ User registration and login
- ✅ JWT-based authentication
- ✅ Protected routes
- ✅ Cross-tab auth synchronization
- ✅ Secure password hashing


## Development

The frontend runs on port 8081 with hot reload. The backend runs on port 8080 with auto-restart using Air (if installed).



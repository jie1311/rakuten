# Rakuten Backend

A Go backend service with user authentication using JWT tokens and MongoDB.

## Prerequisites

- Go 1.24.6+
- MongoDB running locally on port 27017

## Getting Started

### 1. Install Dependencies

```bash
go mod download
```

### 2. Start the Server

```bash
cd backend
go run main.go
```

The server will start on `http://0.0.0.0:8080`

## API Testing with curl

### Sign Up

Create a new user account:

```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Expected response: `201 Created`

### Sign In

Get an authentication token:

```bash
curl -X POST http://localhost:8080/api/auth/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Expected response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "email": "user@example.com"
}
```

### Get User Info (Protected Route)

Use the token from sign-in:

```bash
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

Expected response:
```json
{
  "email": "user@example.com"
}
```

### Sign Out

```bash
curl -X POST http://localhost:8080/api/auth/signout
```

Expected response: `200 OK`

## Running Tests

```bash
cd backend/handler
go test -v
```

## Project Structure

```
├── backend/
│   ├── main.go              # Application entry point
│   ├── database/            # Database interface and MongoDB implementation
│   └── handler/             # HTTP handlers and middleware
└── README.md
```

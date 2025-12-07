# Rakuten Project Report

## Setup and Architectural Choices

### Architecture Overview
The application follows a **client-server architecture** with clear separation of concerns:

**Backend (Go):**
- RESTful API using native `net/http` package
- MongoDB for data persistence
- JWT-based stateless authentication
- Bcrypt for password hashing
- CORS middleware for cross-origin requests

**Frontend (React + TypeScript):**
- Single Page Application (SPA) with React Router
- Context API for state management
- TypeScript for type safety
- Vite for fast development and optimized builds

### Key Architectural Decisions

1. **JWT Authentication**: Chosen for stateless authentication, enabling horizontal scaling without session storage
2. **MongoDB**: Document-oriented database suitable for user data and easy to scale
3. **Context API**: Lightweight state management for authentication state
4. **TypeScript**: Type safety on frontend reduces runtime errors
5. **Modular Structure**: Clear separation between auth logic, routing, and UI components

## Potential Weaknesses and Production Solutions

### Current Weaknesses

1. **No Rate Limiting**
   - *Solution*: Implement rate limiting middleware (e.g., `golang.org/x/time/rate`)
   - Add IP-based throttling for authentication endpoints

2. **JWT Security**
   - No token refresh mechanism
   - *Solution*: Implement refresh tokens with shorter access token expiry
   - Add token blacklist for logout/revocation

3. **No Input Validation**
   - *Solution*: Add validation library (e.g., `go-playground/validator`)
   - Implement email format validation, password strength requirements

4. **Database Connection Management**
   - Single connection, no pooling
   - *Solution*: Implement connection pooling with configurable limits

Weaknesses above is acceptable at this scale as an interview demo, but should be addressed in production with large data

## Future Improvements

If I had more time, I would prioritize:

1. **Authentication Enhancements**
   - Email verification on signup
   - Password reset functionality

2. **Security Hardening**
   - Implement rate limiting
   - Add account lockout after failed attempts
   - Password complexity requirements

3. **Code Quality**
   - Comprehensive test coverage (unit, integration, e2e)
   - API documentation (Swagger/OpenAPI)
   - Code comments and documentation
   - Pre-commit hooks (linting, formatting)

4. **Observability**
   - Structured logging
   - APM integration
   - Error tracking (Sentry)
   - Analytics

## Frontend State and Data-Flow

### State Management Choice: Context API

**Why Context API?**

1. **Simplicity**: For authentication state, Context API is sufficient and doesn't require external dependencies
2. **Built-in**: No additional libraries needed, reducing bundle size
3. **Type-safe**: Works seamlessly with TypeScript
4. **Cross-tab Sync**: Implemented `storage` event listener for synchronized logout across tabs
5. **Right-sized**: The app's state requirements don't justify Redux/MobX complexity

**Alternative Considerations:**
- **Redux**: Overkill for simple auth state
- **Zustand**: Good alternative, but Context API meets current needs
- **React Query**: Would be beneficial for server state management if we had more data fetching

## Types and Contracts

### Current Approach: **Manual Type Management**

**Frontend Types**: Defined manually in TypeScript
```typescript
interface AuthContextType {
  token: string | null
  email: string | null
  setAuth: (token: string, email: string) => void
  clearAuth: () => void
}
```

**Backend Types**: Go structs
```go
type SignupRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

### Production-Ready Approach

**OpenAPI/Swagger**

1. **Generate OpenAPI spec from Go**
   - Use `swaggo/swag` to generate from comments
   - Define schemas, endpoints, request/response types

2. **Generate TypeScript types from OpenAPI**
   - Use `openapi-typescript` to auto-generate TS types
   - Keep frontend types in sync automatically

3. **Benefits**:
   - Single source of truth
   - Automatic type generation
   - API documentation included
   - Contract testing possible

**Alternative Approaches: Shared Type Definitions**
- Use gRPC with Protocol Buffers
- Generate types for both Go and TypeScript
- Strong typing across the stack

## Scenario 1: Brute-Force Attack on Logins

### Current Vulnerability
The login endpoint has no protection against brute-force attacks.

### Mitigation Strategy

#### 1. **Rate Limiting** 
```go
// Implement rate limiter middleware
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Limit(requestsPerMinute), requestsPerMinute*2)
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

#### 2. **Account Lockout** 
- Track failed login attempts per email in database/cache
- Lock account after 5 failed attempts for 15 minutes
- Send email notification on lockout
- Implement exponential backoff

```go
type FailedLogin struct {
    Email      string
    Attempts   int
    LockedUntil time.Time
}
```

#### 3. **IP-Based Rate Limiting** 
- Use Redis to track login attempts per IP
- More aggressive limits per IP address
- Block suspicious IPs temporarily

#### 5. **Monitoring and Alerting**
- Log all failed login attempts
- Alert on suspicious patterns
- Implement honeypot accounts

## Scenario 2: Millions of Requests/Sec & Fault Tolerance

### Architecture Redesign for Scale

#### 1. **Horizontal Scaling**

**Current**: Single server instance
**Solution**: Multiple instances behind load balancer

#### 2. **Caching Strategy**

**Redis for:**
- Session data
- Rate limiting counters
- Hot user data
- API response caching

#### 3. **Database Optimization**

**MongoDB Sharding:**
- Shard by user ID range
- Replica sets for read scaling
- Read preference: secondary for non-critical reads

**Connection Pooling:**
```go
clientOptions := options.Client().
    SetMaxPoolSize(100).
    SetMinPoolSize(10).
    SetMaxConnIdleTime(30 * time.Second)
```

#### 4. **Asynchronous Processing**

**Message Queue (Kafka/RabbitMQ):**
- Email verification
- Password reset emails
- Analytics events
- Audit logging

```go
// Non-blocking operations
go func() {
    queue.Publish("user.signup", SignupEvent{
        Email: user.Email,
        Time:  time.Now(),
    })
}()
```

#### 5. **Auto-scaling**

## Conclusion

This project demonstrates a functional authentication system with room for significant production enhancements. The proposed improvements focus on security, scalability, and operational excellence—critical factors for any production system handling user authentication at scale.

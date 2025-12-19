# Go Gin RealWorld API

A backend implementation of the [RealWorld](https://github.com/gothinkster/realworld) (Conduit) specification using **Go** and **Gin**.

## Tech Stack

- **Language:** Go
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** MySQL
- **Authentication:** JWT (JSON Web Token)
- **Containerization:** Docker & Docker Compose
- **API Documentation:** Swagger (OpenAPI)

## Features

- **Authentication:** User registration, login, and JWT-based authorization.
- **Profiles:** Get user profiles, follow/unfollow users.
- **Articles:** CRUD operations, slug generation, filtering by tag/author/favorited, and personalized feed.
- **Comments:** Add and delete comments on articles.
- **Favorites:** Favorite and unfavorite articles.
- **Tags:** List all unique tags used in articles.

## Error Handling

The project uses a standardized error response format:

```json
{
  "code": 400,
  "message": "Validation failed",
  "details": {
    "email": "must be a valid email",
    "password": "is required"
  }
}
```

- **Validation Errors:** Automatically handled for request binding, providing field-specific error messages.
- **Business Logic Errors:** Defined as constants in `internal/errors/errors.go` for consistency across the application.
- **HTTP Status Codes:** Correctly mapped to error types (e.g., 401 for invalid credentials, 404 for not found, 422/400 for validation).

## Database Transactions

The project ensures data integrity by using database transactions for operations involving multiple steps (e.g., user registration creating both a user and a profile):

- **Service Layer Responsibility:** Transactions are initiated in the service layer using GORM's `Transaction` method.
- **Repository Layer:** Repositories accept a `*gorm.DB` instance (which can be a transaction object) to perform operations within the same unit of work.
- **Automatic Management:** Transactions are automatically committed upon successful completion or rolled back if any error occurs within the block.

**Example (User Registration):**

```go
err := s.db.Transaction(func(tx *gorm.DB) error {
    // 1. Create User
    if err := s.userRepo.CreateUser(tx, user); err != nil {
        return err // Transaction rolls back
    }

    // 2. Create associated Profile
    if err := s.profileRepo.CreateProfile(tx, profile); err != nil {
        return err // Transaction rolls back
    }

    return nil // Transaction commits
})
```

## Getting Started

### Running with Docker

The easiest way to run the project is using Docker Compose:

```bash
docker compose up -d --build
```

The API will be available at `http://localhost:8080`.

### Local Development

1. Clone the repository.
2. Set up a MySQL database.
3. Configure environment variables in a `.env` file (refer to `internal/config/config.go` or `docker-compose.yml`).
4. Run the application:

   ```bash
   go run cmd/app/main.go
   ```

## API Documentation

- **Swagger:** See [swagger.yaml](swagger.yaml) for the API specification.
- **Postman:** A Postman collection is available at [GoGINRealworldAPI.postman_collection.json](GoGINRealworldAPI.postman_collection.json).

## Testing

Run all tests using the following command:

```bash
go test -v ./test/...
```

## Project Structure

The project follows a clean architecture pattern, separating concerns into distinct layers:

```text
.
├── cmd/
│   └── app/
│       └── main.go          # Application entry point
├── internal/
│   ├── bootstrap/           # Dependency injection & container setup
│   ├── config/              # Environment & Database configurations
│   ├── dtos/                # Data Transfer Objects (Request/Response)
│   ├── errors/              # Custom error types & handling logic
│   ├── handlers/            # HTTP controllers (Gin handlers)
│   ├── middleware/          # Gin middlewares (JWT Auth, etc.)
│   ├── models/              # GORM database models
│   ├── repository/          # Data access layer (Interfaces)
│   │   └── mysql/           # MySQL implementations
│   ├── routes/              # API route definitions
│   ├── services/            # Business logic layer
│   └── utils/               # Helper utilities (JWT, Slug, etc.)
├── test/                    # Unit and integration tests
│   └── mocks/               # Mock implementations for testing
├── docker-compose.yml       # Docker orchestration
├── go.mod                   # Go module dependencies
└── swagger.yaml             # OpenAPI 3.0 specification
```

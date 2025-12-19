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

- `cmd/app/`: Application entry point.
- `internal/`: Core business logic, handlers, models, and repositories.
- `test/`: Unit and integration tests.
- `docker-compose.yml`: Docker configuration for the app and database.

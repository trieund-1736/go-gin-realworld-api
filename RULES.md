# Folder Structure

## Generic Structure (Template)

```text
my-go-project/
│
├── cmd/
│   └── app/
│       └── main.go              # Entry point (start server, router, config)
│
├── go.mod                       # Go modules (dependencies)
├── go.sum
├── README.md
│
├── internal/                    # Code ONLY used for this project
│   ├── config/                  # Load & parse config
│   │   ├── config.go
│   │   └── database.go
│   │
│   ├── handlers/                # HTTP handlers / controllers
│   │   ├── user_handler.go
│   │   └── product_handler.go
│   │
│   ├── middleware/              # Middleware (auth, logging, recover...)
│   │   ├── auth.go
│   │   └── logger.go
│   │
│   ├── repository/              # DB access layer (CRUD, query)
│   │   ├── user_repository.go
│   │   └── product_repository.go
│   │
│   ├── models/                  # Data models / entities
│   │   ├── user.go
│   │   └── product.go
│   │
│   ├── services/                # Business logic (optional but highly recommended)
│   │   ├── user_service.go
│   │   └── product_service.go
│   │
│   └── routes/                  # Route definitions
│       ├── user_routes.go
│       └── product_routes.go
│
├── pkg/                         # Package REUSABLE for other projects
│   ├── logger/
│   │   └── logger.go
│   └── validator/
│       └── validator.go
│
├── config/                      # Configuration files
│   ├── app.yaml
│   └── database.yaml
│
├── env/                         # Environment-specific config
│   ├── dev.env
│   ├── staging.env
│   └── prod.env
│
├── utils/                       # Common utility functions
│   ├── time.go
│   └── string.go
│
├── helpers/                     # Helpers for specific workflows
│   └── response_helper.go
│
├── static/                      # Static assets
│   ├── css/
│   ├── js/
│   └── images/
│
├── templates/                   # HTML templates
│   ├── layout.html
│   └── index.html
│
└── tests/                       # Unit / integration tests
    └── user_test.go

```

## Current Project Structure (go-gin-realworld-api)

```text
go-gin-realworld-api/
│
├── cmd/
│   └── app/
│       └── main.go                      # Entry point
│
├── internal/
│   ├── bootstrap/
│   │   └── container.go                 # DI container setup
│   │
│   ├── config/
│   │   ├── config.go                    # Configuration management
│   │   └── database.go                  # Database connection & migration
│   │
│   ├── dtos/                            # Request/Response DTOs for API
│   │
│   ├── handlers/                        # HTTP handlers/controllers for routes
│   │
│   ├── middleware/
│   │   └── auth.go                      # JWT authentication
│   │
│   ├── models/                          # GORM models (match DATABASE.md schema)
│   │
│   ├── repository/                      # Data access layer (CRUD operations)
│   │
│   ├── services/                        # Business logic layer
│   │
│   ├── routes/
│   │   └── routes.go                    # Route definitions & registration
│   │
│   └── utils/
│       └── jwt.go                       # JWT utilities
│
├── pkg/
│   └── jwt/                             # Reusable JWT package
│
├── .env.example                         # Environment variables template
├── DATABASE.md                          # Database schema
├── REQUIREMENT.md                       # API requirements
├── RULES.md                             # This file
├── swagger.yaml                         # OpenAPI specification
└── go.mod, go.sum                       # Go dependencies
```

---

## Workflow: Adding New API Features

1. **Database Schema** → Add table to [DATABASE.md](DATABASE.md)
2. **Model** → Create `{feature}.go` in [internal/models/](internal/models) (see style from [user.go](internal/models/user.go), [article.go](internal/models/article.go))
3. **Migration** → Add model to `AutoMigrate()` in [internal/config/database.go](internal/config/database.go)
4. **DTO** → Create `{feature}_dto.go` in [internal/dtos/](internal/dtos) (see [user_dto.go](internal/dtos/user_dto.go), [auth_dto.go](internal/dtos/auth_dto.go))
5. **Repository** → Create `{feature}_repository.go` in [internal/repository/](internal/repository) (see [user_repository.go](internal/repository/user_repository.go))
6. **Service** → Create `{feature}_service.go` in [internal/services/](internal/services) (see [auth_service.go](internal/services/auth_service.go))
7. **Handler** → Create `{feature}_handler.go` in [internal/handlers/](internal/handlers) (see [auth_handler.go](internal/handlers/auth_handler.go))
8. **Routes** → Register routes in [internal/routes/routes.go](internal/routes/routes.go)
9. **Swagger** → Update endpoints in [swagger.yaml](swagger.yaml)
10. **Bootstrap** → Register handler, service, repository in [internal/bootstrap/container.go](internal/bootstrap/container.go)

---

## Key Rules

- **Models**: GORM struct with relationships (foreignKey, constraint, preload)
- **DTOs**: Separate request & response types, use `binding` tags for validation
- **Repositories**: CRUD only, no business logic
- **Services**: Business logic, validation, error handling
- **Handlers**: HTTP logic, binding request, format response
- **Routes**: Group by feature, apply middleware (auth) if needed
- **Swagger**: Update when adding/modifying endpoint

# Folder Structure

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

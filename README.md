# Go Framework

An enterprise-grade, modular Go framework designed for building scalable web applications.

## Features

- **Modular Design**: Decoupled components for HTTP, Database, Logging, Config, and CLI.
- **HTTP Server**: Built on Gin with custom middleware support.
- **Database**: Integrated Gorm with repository pattern support (MySQL, SQLite).
- **Logging**: Structured logging with Zap and log rotation (Lumberjack).
- **Configuration**: Viper-based configuration management (YAML, JSON, Env, Flags).
- **CLI**: Cobra-based CLI for project scaffolding and management.

## Installation

```bash
go install github.com/project/go-framework/cmd/goframework@latest
```

## Quick Start

1. **Create a new project**:
   ```bash
   goframework new my-app
   ```

2. **Run the project**:
   ```bash
   cd my-app
   go mod tidy
   go run cmd/server/main.go
   ```

## Project Structure

```
my-app/
├── cmd/
│   └── server/          # Main application entry point
├── configs/             # Configuration files
├── internal/
│   ├── api/             # HTTP handlers
│   ├── models/          # Database models
│   └── service/         # Business logic
├── pkg/                 # Public libraries
├── go.mod
└── Makefile
```

## Documentation

### Configuration
The framework uses `config.yaml` by default. Example:
```yaml
app:
  name: my-app
  version: 1.0.0
server:
  port: 8080
  mode: debug
log:
  level: info
database:
  driver: sqlite
  source: test.db
```

### Database Access
Use the generic `Repository` for CRUD operations:
```go
repo := database.NewRepository[models.User](db)
err := repo.Create(&user)
user, err := repo.FindByID(1)
```

## License
MIT

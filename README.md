# Relay API

A fast, lightweight REST API server written in Go with the Fiber framework.

## Features

- JWT-based authentication
- RESTful CRUD endpoints
- PostgreSQL integration
- Structured JSON logging
- Graceful shutdown

## Quick Start

```bash
go mod download
go run .
```

Server starts on [http://localhost:8080](http://localhost:8080).

## Endpoints

| Method | Path           | Description         | Auth  |
|--------|----------------|---------------------|-------|
| POST   | /api/auth/login | Get JWT token       | No    |
| GET    | /api/health     | Health check        | No    |
| GET    | /api/items      | List all items      | Yes   |
| POST   | /api/items      | Create an item      | Yes   |
| GET    | /api/items/:id  | Get item by ID      | Yes   |
| DELETE | /api/items/:id  | Delete item by ID   | Yes   |

## Environment Variables

| Variable     | Default           | Description       |
|-------------|-------------------|-------------------|
| PORT        | 8080              | Server port       |
| DB_URL      | postgres://...    | PostgreSQL DSN    |
| JWT_SECRET  | (required)        | JWT signing key   |

## License

MIT

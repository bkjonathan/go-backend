# Go Backend (Production-Ready Scaffold)

## Stack
- Router: chi
- DB access: GORM + PostgreSQL
- Config: cleanenv
- Validation: validator/v10
- Logging: zap
- Auth: JWT + bcrypt
- Runtime: Docker + Docker Compose

## Quick Start
1. Copy `.env.example` to `.env`.
2. Run services:
    - `make docker-up`
3. Start app locally:
    - `make run`

The app runs `AutoMigrate` on startup and manages the `users` table automatically.

## API
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/users` (protected)
- `GET /api/v1/users/me` (protected)
- `GET /api/v1/users/{id}` (protected)
- `PUT /api/v1/users/{id}` (protected)
- `DELETE /api/v1/users/{id}` (protected)

# Imphnen Backend API

A Go-based backend service for managing users, products, orders, and transactions with a Telegram bot integration.

## Features

- RESTful API for user management
- Product and order management
- Transaction processing and analytics
- Telegram bot integration for customer operations
- JWT authentication
- Cloudinary integration for image storage
- Swagger documentation

## Requirements

- Go 1.25+
- PostgreSQL 16+
- Docker (optional)

## Setup

1. Copy `.env.example` to `.env` and fill in the required environment variables:

   ```bash
   cp .env.example .env
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run the application:

   ```bash
   go run cmd/main.go
   ```

## API Documentation

Access the Swagger UI at `http://localhost:8080/api/v1/docs/`

## Docker Deployment

1. Build and run the application:

   ```bash
   docker-compose up --build
   ```

2. The application will be available at `http://localhost:6969`

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DSN`: PostgreSQL database connection string
- `JWT_SECRET`: JWT signing secret
- `KOLOSAL_API_KEY`: Kolosal API key
- `CLOUDINARY_NAME`: Cloudinary account name
- `CLOUDINARY_API_KEY`: Cloudinary API key
- `CLOUDINARY_API_SECRET`: Cloudinary API secret

## License

MIT
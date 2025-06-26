# URL Shortener Service

A simple URL shortener service built with Gin in Go.

## Features

- Generate short URLs from long URLs
- Redirect short URLs to original long URLs
- Handle duplicate URLs by returning existing short codes
- Comprehensive unit tests
- Docker containerization on mysql
- Support for .env files

## Project Structure

```
.
├── cmd/           # Configuration management
│   └── shortener/
│       └── main.go    # Viper-based config loading
├── internal/          # Business logic
│   ├── app/              # Application layer: Contains core business logic (e.g., ShortenerService).
│   │   ├── errors.go
│   │   ├── shortener_service.go
│   │   └── shortener_service_test.go
│   ├── domain/           # Domain layer: Defines core entities and interfaces (e.g., ShortURL, ShortURLRepository).
│   │   ├── short_url.go
│   │   └── short_url_repository.go
│   └── infrastructure/   # Infrastructure layer: Contains external dependencies and utilities.
│       ├── id_generator/ # Base62 encoding for short codes.
│       ├── persistence/  # Database implementations (e.g., MySQL GORM repository).
│       │   └── mysql/
│       │       └── mysql_short_url_repository.go
│       └── web/          # Web layer: HTTP handlers and router setup (Gin Gonic).
│           ├── handler.go
│           └── router.go
├── docker-compose.yaml        # Container configuration
│
├── go.mod            # Go module definition
├── go.sum            # Dependency checksums
└── README.md         # This file
```

### Environment Variables

| Variable      | Default                 | Description                              |
| ------------- | ----------------------- | ---------------------------------------- |
| `PORT`        | `8080`                  | Port number for the server               |
| `BASE_URL`    | `http://localhost:8080` | Base URL for generating short links      |
| `LOG_LEVEL`   | `info`                  | Logging level (debug, info, warn, error) |

### Example .env File

Create a `.env` file in the project root:

```bash
# Dev settings
DB_PASSWORD=root
DB_ROOT_PASSWORD=root
DB_HOST=127.0.0.1
DB_PORT=33066
DB_NAME=shortener_db
```

## API Endpoints

### 1. Shorten URL

**POST** `/api/shortlinks`

Request Body:

```json
{
  "long_url": "https://www.example.com/very/long/url"
}
```

Response (New URL - 201 Created):

```json
{
  "id": "abc123"
}
```

Response (Duplicate URL - 200 OK):

```json
{
  "id": "abc123"
}
```

### 2. Get Short Link Details

**GET** `/api/shortlinks/{shortId}`

Response:

```json
{
  "id": "abc123"
}
```

### 3. Redirect to Long URL

**GET** `/shortlinks/{id}`

Public redirect endpoint – 302 redirect to the original URL

### Prerequisites

- Go 1.24.1 or later
- Docker and Docker compose (for containerized deployment)

### Initial Setup

```bash
# Setup the project (install dependencies)
docker compose up -d (for mysql only)
```

## Notes

- Duplicate URLs are handled efficiently - the same URL will always return the same short code.
- The Docker image for dev only, not for production

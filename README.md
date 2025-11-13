# Student API

A RESTful API for managing student records, built with Go using SQLite for persistence.

## Features

- ✅ Create new students
- ✅ Retrieve all students or a specific student by ID
- ✅ Update student information (partial updates)
- ✅ Delete students
- ✅ Input validation with detailed error messages
- ✅ Structured JSON logging
- ✅ Graceful server shutdown
- ✅ Configuration management via YAML and environment variables

## Prerequisites

- Go 1.25.3 or higher
- SQLite3

## Project Structure

```
student-api/
├── cmd/
│   └── student-api/
│       └── main.go              # Application entry point
├── config/
│   └── local.yml                # Configuration file
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading
│   ├── http/
│   │   └── handlers/
│   │       └── student/
│   │           └── student.go   # HTTP handlers
│   ├── storage/
│   │   ├── storage.go           # Storage interface
│   │   └── sqlite/
│   │       └── sqlite.go        # SQLite implementation
│   ├── types/
│   │   └── types.go             # Data structures
│   └── utils/
│       └── response/
│           └── response.go      # Response utilities
├── storage/                     # SQLite database file (created at runtime)
├── go.mod                       # Go module dependencies
└── .gitignore                   # Git ignore rules
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/gourav224/student-api.git
cd student-api
```

2. Install dependencies:
```bash
go mod download
```

## Configuration

Configuration can be loaded from:
1. Environment variable: `CONFIG_PATH`
2. Command-line flag: `--config`

Example configuration in `config/local.yml`:
```yaml
env: "dev"
storage_path: "storage/sqlite.db"
http_server:
  address: "localhost:8000"
```

### Environment Variables

- `CONFIG_PATH`: Path to the configuration file
- `HTTP_SERVER_ADDR`: HTTP server address (default: `:8080`)
- `STORAGE_PATH`: SQLite database file path (required)
- `ENV`: Environment name (required)

## Running the Application

```bash
CONFIG_PATH=config/local.yml go run cmd/student-api/main.go
```

Or with the config flag:
```bash
go run cmd/student-api/main.go --config=config/local.yml
```

The server will start on `localhost:8000` (as per `config/local.yml`).

## API Endpoints

### Create Student
**POST** `/api/students`

Request body:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 20
}
```

Response (201 Created):
```json
{
  "status": "success",
  "message": "student created successfully",
  "data": 1
}
```

### Get All Students
**GET** `/api/students`

Response (200 OK):
```json
{
  "status": "success",
  "message": "students fetched successfully",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "age": 20
    }
  ]
}
```

### Get Student by ID
**GET** `/api/students/{id}`

Response (200 OK):
```json
{
  "status": "success",
  "message": "student fetched successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 20
  }
}
```

### Update Student
**PATCH** `/api/students/{id}`

Request body (all fields optional):
```json
{
  "name": "Jane Doe",
  "age": 21
}
```

Response (200 OK):
```json
{
  "status": "success",
  "message": "student updated successfully",
  "data": {
    "id": 1,
    "name": "Jane Doe",
    "email": "john@example.com",
    "age": 21
  }
}
```

### Delete Student
**DELETE** `/api/students/{id}`

Response (200 OK):
```json
{
  "status": "success",
  "message": "student deleted successfully",
  "data": 1
}
```

## Validation Rules

The following validation rules are enforced:

- **name**: Required, must be a non-empty string
- **email**: Required, must be a valid email address
- **age**: Required, must be an integer between 1 and 120

## Dependencies

- `github.com/go-playground/validator/v10` - Input validation
- `github.com/ilyakaznacheev/cleanenv` - Configuration management
- `github.com/mattn/go-sqlite3` - SQLite3 driver

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful GET/PATCH/DELETE
- `201 Created` - Successful POST
- `400 Bad Request` - Invalid input or malformed request
- `500 Internal Server Error` - Database or server errors

All error responses follow this format:
```json
{
  "status": "error",
  "error": "error message"
}
```

## Development

### Building

```bash
go build -o student-api cmd/student-api/main.go
```

### Running Tests

```bash
go test ./...
```

## License

This project is open source and available under the MIT License.

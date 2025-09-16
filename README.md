# BAGR Backend System

A clean, modular Go backend for a bidding/music marketplace application built with modern practices and scalable architecture.

## 🏗️ Architecture

This project follows a clean architecture pattern with clear separation of concerns:

```
├── cmd/                    # Application entrypoints
│   └── main.go            # Main application entry
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── controllers/      # HTTP request handlers
│   ├── models/          # Data models and DTOs
│   ├── repositories/    # Data access layer
│   ├── server/          # HTTP server setup and routing
│   ├── services/        # Business logic layer
│   └── utils/           # Utility functions and helpers
├── Dockerfile           # Container configuration
├── Makefile            # Development and build commands
├── config.yaml         # Configuration file
└── env.example         # Environment variables example
```

## 🚀 Features

- **Clean Architecture**: Modular design with clear separation of concerns
- **HTTP API**: RESTful API built with Gin framework
- **Configuration Management**: Environment-based configuration with YAML support
- **Structured Logging**: JSON logging with Logrus
- **Middleware**: CORS, logging, recovery, request ID, and timeout middleware
- **Health Checks**: Health and readiness endpoints
- **Docker Support**: Multi-stage Docker build for production
- **Development Tools**: Comprehensive Makefile with common tasks

## 🛠️ Tech Stack

- **Go 1.21+**: Modern Go with latest features
- **Gin**: Fast HTTP web framework
- **PostgreSQL**: Primary database (ready for integration)
- **Redis**: Caching and session storage (ready for integration)
- **Logrus**: Structured logging
- **Docker**: Containerization

## 📋 Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- PostgreSQL (for database features)
- Redis (for caching features)

## 🏃‍♂️ Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd BAGR_Backend_System
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Configure Environment

Copy the example environment file and modify as needed:

```bash
cp env.example .env
# Edit .env with your configuration
```

### 4. Run the Application

```bash
# Run with default configuration
make run

# Or run with custom config file
make run-with-config
```

The server will start on `http://localhost:8080`

## 🐳 Docker Deployment

### Build and Run with Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

### Using Docker Compose (Future)

```bash
# Start all services (when docker-compose.yml is added)
make db-up
make docker-run
```

## 📚 API Documentation

### Health Endpoints

- `GET /health` - Health check
- `GET /ready` - Readiness check

### User Endpoints

- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users (with pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Example API Calls

```bash
# Health check
curl http://localhost:8080/health

# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "password": "securepassword",
    "role": "buyer"
  }'

# Get user by ID
curl http://localhost:8080/api/v1/users/1

# List users with pagination
curl "http://localhost:8080/api/v1/users?limit=10&offset=0"
```

## 🧪 Development

### Available Make Commands

```bash
make help                 # Show all available commands
make run                  # Run the application
make build               # Build binary
make test                # Run tests
make test-coverage       # Run tests with coverage
make fmt                 # Format code
make vet                 # Run go vet
make lint                # Run golangci-lint
make dev-check           # Run all development checks
make docker-build        # Build Docker image
make docker-run          # Run Docker container
```

### Development Workflow

```bash
# Setup development environment
make dev-setup

# Run development checks and start server
make dev-run

# Or run checks individually
make fmt vet lint test
```

## 📁 Project Structure Details

### Models
- `User`: User account management
- `Auction`: Auction listings
- `Bid`: Bidding system
- `Track`: Music track metadata

### Services
Business logic layer that handles:
- User management
- Authentication (ready for implementation)
- Auction management (ready for implementation)
- Bidding logic (ready for implementation)

### Repositories
Data access layer with interfaces for:
- Database operations
- Caching operations
- External API integrations

### Controllers
HTTP request handlers that:
- Validate input
- Call appropriate services
- Return structured responses

## 🔧 Configuration

The application supports configuration through:

1. **YAML file** (`config.yaml`)
2. **Environment variables** (`.env` file)
3. **Command line flags**

Environment variables take precedence over YAML configuration.

### Configuration Options

- **Server**: Host, port, timeouts
- **Database**: PostgreSQL connection settings
- **Redis**: Cache configuration
- **Application**: Environment, logging, JWT secret

## 🚦 Future Enhancements

This project is designed to be easily extensible. Planned features include:

- **Authentication & Authorization**: JWT-based auth system
- **Auction System**: Complete bidding functionality
- **Music Library**: Track upload and management
- **Auto-bidder**: Automated bidding system
- **WebSocket Support**: Real-time bidding updates
- **File Upload**: Music file and image handling
- **Payment Integration**: Payment processing
- **Email Notifications**: User notifications
- **Admin Dashboard**: Management interface

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the API examples above

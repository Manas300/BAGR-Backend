# 🎵 BAGR Backend - Beat Auction & Global Records

A comprehensive Go-based backend system for a music auction platform with real-time bidding, user authentication, and beat streaming capabilities.

## 🚀 Features

### 🔐 Authentication System
- **JWT-based authentication** with access and refresh tokens
- **Role-based access control** (Admin, Moderator, Producer, Artist, Fan)
- **Email verification** with HTML templates
- **Password reset** functionality
- **Secure password validation** with strength requirements
- **Comprehensive logging** for debugging and monitoring

### 🎵 Music Auction System
- **Real-time bidding** with WebSocket support
- **Beat streaming** capabilities
- **Auction management** with time-based events
- **User profiles** and track management
- **Bid tracking** and history

### 🛠️ Technical Features
- **Clean Architecture** with separation of concerns
- **PostgreSQL** database with migrations
- **Redis** for caching and real-time features
- **Docker** containerization
- **Comprehensive logging** with structured JSON logs
- **Email service** with SMTP support
- **RESTful API** with proper error handling

## 📁 Project Structure

```
BAGR-Backend/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── auth/                   # Authentication package
│   │   ├── handlers.go         # HTTP handlers for auth endpoints
│   │   ├── service.go          # Authentication business logic
│   │   ├── jwt.go             # JWT token management
│   │   ├── password.go        # Password validation and hashing
│   │   └── email.go           # Email service with SMTP
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── models/
│   │   ├── user.go            # User model and DTOs
│   │   ├── auction.go         # Auction model
│   │   ├── bid.go             # Bid model
│   │   └── track.go           # Track model
│   ├── repositories/
│   │   ├── interfaces.go      # Repository interfaces
│   │   └── user_repository.go # User data access layer
│   ├── services/
│   │   └── user_service.go    # User business logic
│   ├── server/
│   │   ├── server.go          # HTTP server setup
│   │   ├── routes.go          # API route definitions
│   │   └── middleware.go      # Custom middleware
│   └── utils/
│       ├── logger.go          # Logging utilities
│       └── response.go        # HTTP response helpers
├── migrations/
│   ├── 000_base_tables.sql    # Base database schema
│   └── 001_auth_tables.sql    # Authentication tables
├── templates/
│   └── verification_success.html # Email verification template
├── config.yaml                # Application configuration
├── docker-compose.yml         # Docker services
├── Dockerfile                 # Container configuration
└── go.mod                     # Go module dependencies
```

## 🛠️ Installation & Setup

### Prerequisites
- **Go 1.21+**
- **PostgreSQL 13+**
- **Redis 6+**
- **Docker** (optional)

### 1. Clone the Repository
```bash
git clone https://github.com/Manas300/BAGR-Backend.git
cd BAGR-Backend
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Database Setup
```bash
# Start PostgreSQL and Redis with Docker
docker-compose up -d

# Run database migrations
psql -h localhost -U bagr_user -d bagr_db -f migrations/000_base_tables.sql
psql -h localhost -U bagr_user -d bagr_db -f migrations/001_auth_tables.sql
```

### 4. Configuration
Update `config.yaml` with your settings:
```yaml
server:
  host: "localhost"
  port: 8080

database:
  host: "localhost"
  port: 5432
  name: "bagr_db"
  username: "bagr_user"
  password: "bagr_password"

jwt:
  access_secret: "your-access-secret-key"
  refresh_secret: "your-refresh-secret-key"

email:
  host: "smtp.gmail.com"
  port: "587"
  username: "your-email@gmail.com"
  password: "your-app-password"
  from_email: "your-email@gmail.com"
  from_name: "BAGR Auction System"
  test_mode: false
```

### 5. Run the Application
```bash
go run cmd/main.go -config config.yaml
```

## 📚 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/verify` - Email verification
- `POST /api/v1/auth/forgot-password` - Password reset request
- `POST /api/v1/auth/reset-password` - Password reset
- `POST /api/v1/auth/refresh` - Token refresh
- `GET /api/v1/auth/profile` - Get user profile (protected)
- `PUT /api/v1/auth/profile` - Update user profile (protected)
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/roles` - Get available roles

### Health Check
- `GET /health` - Health check endpoint
- `GET /ready` - Readiness check endpoint

## 🔧 User Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| **Admin** | Platform administrators | Full system access |
| **Moderator** | Platform moderators | Content moderation, user management |
| **Producer** | Music creators who sell beats | Create auctions, manage tracks |
| **Artist** | Music creators who buy beats | Bid on auctions, purchase tracks |
| **Fan** | General users | View auctions, limited participation |

## 🔐 Authentication Flow

1. **Registration**: User provides email, username, password, and role
2. **Email Verification**: System sends verification email with token
3. **Login**: User logs in with verified credentials
4. **JWT Tokens**: System issues access and refresh tokens
5. **Protected Routes**: Access with valid JWT token

## 📧 Email System

The system supports both test mode and production email sending:

- **Test Mode**: Logs email content to console (for development)
- **Production Mode**: Sends actual emails via SMTP

## 🐳 Docker Support

```bash
# Build and run with Docker Compose
docker-compose up --build

# Run in background
docker-compose up -d
```

## 🧪 Testing

The project includes comprehensive test suites:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
cd tests && ./run_all_tests.ps1
```

## 📊 Logging

The application uses structured JSON logging with different levels:

- **INFO**: General application flow
- **DEBUG**: Detailed debugging information
- **ERROR**: Error conditions
- **WARN**: Warning conditions

## 🔧 Development

### Code Structure
- **Clean Architecture** with clear separation of concerns
- **Dependency Injection** for testability
- **Interface-based design** for flexibility
- **Comprehensive error handling** throughout

### Adding New Features
1. Define models in `internal/models/`
2. Create repository interfaces in `internal/repositories/`
3. Implement business logic in `internal/services/`
4. Add HTTP handlers in `internal/server/`
5. Update routes in `internal/server/routes.go`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Manas Singh** - *Initial work* - [Manas300](https://github.com/Manas300)

## 🙏 Acknowledgments

- Go community for excellent libraries
- Gin framework for HTTP routing
- PostgreSQL and Redis for data storage
- Docker for containerization

---

**Built with ❤️ for the music community**
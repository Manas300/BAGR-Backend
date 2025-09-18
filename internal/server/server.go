package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"bagr-backend/internal/auth"
	"bagr-backend/internal/config"
	"bagr-backend/internal/repositories"
	"bagr-backend/internal/services"
	"bagr-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
)

// Server represents the HTTP server
type Server struct {
	config     *config.Config
	httpServer *http.Server
	db         *sql.DB
}

// Services holds all service instances
type Services struct {
	User    *services.UserService
	Auth    *auth.AuthService
	Profile *services.ProfileService
	S3      *services.S3Service
	Logger  *logrus.Logger
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Initialize logger
	utils.InitLogger(s.config.App.LogLevel)
	logger := utils.GetLogger()

	// Set Gin mode based on environment
	if s.config.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	if err := s.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	repos := s.initRepositories()

	// Initialize services
	services := s.initServices(repos)

	// Initialize controllers
	controllers := NewControllers(services)

	// Create Gin router
	router := gin.New()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Add middleware
	router.Use(LoggerMiddleware())
	router.Use(RecoveryMiddleware())
	router.Use(CORSMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(TimeoutMiddleware(30 * time.Second))

	// Add JWT service to context for middleware
	router.Use(func(c *gin.Context) {
		c.Set("jwt_service", services.Auth.GetJWTService())
		c.Next()
	})

	// Setup routes
	SetupRoutes(router, controllers)

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         s.config.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	logger.WithField("address", s.config.GetServerAddr()).Info("Starting HTTP server")

	// Start server
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	logger := utils.GetLogger()
	logger.Info("Shutting down HTTP server")

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Failed to shutdown HTTP server gracefully")
		return err
	}

	// Close database connection
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close database connection")
			return err
		}
	}

	logger.Info("Server shutdown completed")
	return nil
}

// initDatabase initializes the database connection
func (s *Server) initDatabase() error {
	logger := utils.GetLogger()

	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", s.config.GetDatabaseURL())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	s.db = db
	logger.Info("Database connection established")

	return nil
}

// initRepositories initializes all repositories
func (s *Server) initRepositories() *repositories.Repositories {
	return &repositories.Repositories{
		User: repositories.NewUserRepository(s.db),
		// Add other repositories here when implemented
	}
}

// initServices initializes all services
func (s *Server) initServices(repos *repositories.Repositories) *Services {
	// Initialize logger
	logger := utils.GetLogger()

	// Initialize auth services
	jwtService := auth.NewJWTService(s.config.JWT.AccessSecret, s.config.JWT.RefreshSecret)
	passwordService := auth.NewPasswordService()
	emailService := auth.NewEmailService(auth.EmailConfig{
		ClientID:     s.config.Email.ClientID,
		ClientSecret: s.config.Email.ClientSecret,
		TenantID:     s.config.Email.TenantID,
		FromEmail:    s.config.Email.FromEmail,
		FromName:     s.config.Email.FromName,
		TestMode:     s.config.Email.TestMode, // Use config value
	})
	authService := auth.NewAuthService(s.db, jwtService, passwordService, emailService)

	// Initialize S3 service
	s3Service, err := services.NewS3Service(
		s.config.S3.Region,
		s.config.S3.Bucket,
		s.config.S3.AccessKeyID,
		s.config.S3.SecretAccessKey,
		s.config.S3.BaseURL,
		logger,
	)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize S3 service")
	}

	// Initialize profile service
	profileService := services.NewProfileService(s.db, logger)

	return &Services{
		User:    services.NewUserService(repos.User),
		Auth:    authService,
		Profile: profileService,
		S3:      s3Service,
		Logger:  logger,
	}
}

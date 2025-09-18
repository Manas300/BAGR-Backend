package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	App      AppConfig      `yaml:"app"`
	JWT      JWTConfig      `yaml:"jwt"`
	Email    EmailConfig    `yaml:"email"`
	S3       S3Config       `yaml:"s3"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string `yaml:"host" env:"SERVER_HOST"`
	Port         string `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST"`
	Port     string `yaml:"port" env:"REDIS_PORT"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Environment string `yaml:"environment" env:"APP_ENV"`
	LogLevel    string `yaml:"log_level" env:"LOG_LEVEL"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	AccessSecret  string `yaml:"access_secret" env:"JWT_ACCESS_SECRET"`
	RefreshSecret string `yaml:"refresh_secret" env:"JWT_REFRESH_SECRET"`
}

// EmailConfig holds email configuration
type EmailConfig struct {
	ClientID     string `yaml:"client_id" env:"EMAIL_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" env:"EMAIL_CLIENT_SECRET"`
	TenantID     string `yaml:"tenant_id" env:"EMAIL_TENANT_ID"`
	FromEmail    string `yaml:"from_email" env:"EMAIL_FROM_EMAIL"`
	FromName     string `yaml:"from_name" env:"EMAIL_FROM_NAME"`
	TestMode     bool   `yaml:"test_mode" env:"EMAIL_TEST_MODE"`
}

// S3Config holds AWS S3 configuration
type S3Config struct {
	Region          string `yaml:"region" env:"S3_REGION"`
	Bucket          string `yaml:"bucket" env:"S3_BUCKET"`
	AccessKeyID     string `yaml:"access_key_id" env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `yaml:"secret_access_key" env:"S3_SECRET_ACCESS_KEY"`
	BaseURL         string `yaml:"base_url" env:"S3_BASE_URL"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{}

	// Load from YAML file if provided
	if configPath != "" {
		if err := loadFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Override with environment variables
	loadFromEnv(config)

	// Set defaults
	setDefaults(config)

	return config, nil
}

// loadFromFile loads configuration from YAML file
func loadFromFile(config *Config, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(config)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *Config) {
	// Server config
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if timeout := os.Getenv("SERVER_READ_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			config.Server.ReadTimeout = val
		}
	}
	if timeout := os.Getenv("SERVER_WRITE_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			config.Server.WriteTimeout = val
		}
	}

	// Database config
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// Redis config
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		config.Redis.Port = port
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if val, err := strconv.Atoi(db); err == nil {
			config.Redis.DB = val
		}
	}

	// App config
	if env := os.Getenv("APP_ENV"); env != "" {
		config.App.Environment = env
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.App.LogLevel = logLevel
	}

	// JWT config
	if accessSecret := os.Getenv("JWT_ACCESS_SECRET"); accessSecret != "" {
		config.JWT.AccessSecret = accessSecret
	}
	if refreshSecret := os.Getenv("JWT_REFRESH_SECRET"); refreshSecret != "" {
		config.JWT.RefreshSecret = refreshSecret
	}

	// Email config
	if clientID := os.Getenv("EMAIL_CLIENT_ID"); clientID != "" {
		config.Email.ClientID = clientID
	}
	if clientSecret := os.Getenv("EMAIL_CLIENT_SECRET"); clientSecret != "" {
		config.Email.ClientSecret = clientSecret
	}
	if tenantID := os.Getenv("EMAIL_TENANT_ID"); tenantID != "" {
		config.Email.TenantID = tenantID
	}
	if fromEmail := os.Getenv("EMAIL_FROM_EMAIL"); fromEmail != "" {
		config.Email.FromEmail = fromEmail
	}
	if fromName := os.Getenv("EMAIL_FROM_NAME"); fromName != "" {
		config.Email.FromName = fromName
	}
	if testMode := os.Getenv("EMAIL_TEST_MODE"); testMode != "" {
		if val, err := strconv.ParseBool(testMode); err == nil {
			config.Email.TestMode = val
		}
	}

	// S3 config
	if region := os.Getenv("S3_REGION"); region != "" {
		config.S3.Region = region
	}
	if bucket := os.Getenv("S3_BUCKET"); bucket != "" {
		config.S3.Bucket = bucket
	}
	if accessKeyID := os.Getenv("S3_ACCESS_KEY_ID"); accessKeyID != "" {
		config.S3.AccessKeyID = accessKeyID
	}
	if secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY"); secretAccessKey != "" {
		config.S3.SecretAccessKey = secretAccessKey
	}
	if baseURL := os.Getenv("S3_BASE_URL"); baseURL != "" {
		config.S3.BaseURL = baseURL
	}
}

// setDefaults sets default values for configuration
func setDefaults(config *Config) {
	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30
	}

	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == "" {
		config.Database.Port = "5432"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	if config.Redis.Host == "" {
		config.Redis.Host = "localhost"
	}
	if config.Redis.Port == "" {
		config.Redis.Port = "6379"
	}

	if config.App.Environment == "" {
		config.App.Environment = "development"
	}
	if config.App.LogLevel == "" {
		config.App.LogLevel = "info"
	}

	// JWT defaults
	if config.JWT.AccessSecret == "" {
		config.JWT.AccessSecret = "your-access-secret-key-change-in-production"
	}
	if config.JWT.RefreshSecret == "" {
		config.JWT.RefreshSecret = "your-refresh-secret-key-change-in-production"
	}

	// Email defaults
	if config.Email.FromEmail == "" {
		config.Email.FromEmail = "admin@bagr.app"
	}
	if config.Email.FromName == "" {
		config.Email.FromName = "BAGR Auction System"
	}
	// TestMode defaults to false (real email sending)
	// Only set to true if explicitly configured

	// S3 defaults
	if config.S3.Region == "" {
		config.S3.Region = "us-east-1"
	}
	if config.S3.Bucket == "" {
		config.S3.Bucket = "bagr-profile-images"
	}
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// GetServerAddr returns the server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
